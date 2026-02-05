package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"photobooth-core/internal/platform/auth"
	"photobooth-core/internal/platform/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware jadi pintu masuk utama buat validasi JWT.
// Di sini kita juga sekalian beresin tipe data ID dari string ke UUID
// supaya di level handler kita nggak perlu repot parsing lagi.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Otorisasi diperlukan", nil)
			c.Abort()
			return
		}

		// Pastikan format header pake "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "Format header salah", nil)
			c.Abort()
			return
		}

		// Validasi token pake secret key yang ada di platform/auth
		token, err := auth.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Token expired atau nggak valid", err.Error())
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Gagal baca payload token", nil)
			c.Abort()
			return
		}

		// Ambil role buat bedain akses (default ke user kalau kosong)
		role, _ := claims["role"].(string)
		if role == "" {
			role = "user"
		}
		c.Set("role", role)

		// Parse Tenant ID ke UUID asli supaya handler nggak panic pas type assertion
		if tIDStr, ok := claims["tenant_id"].(string); ok {
			tenantID, _ := uuid.Parse(tIDStr)
			c.Set("tenant_id", tenantID)
		}

		// Pisahin ID yang disimpen di context berdasarkan rolenya
		if role == "device" {
			// Kalau mesin, kita simpen booth_id-nya
			if bIDStr, ok := claims["booth_id"].(string); ok {
				boothID, _ := uuid.Parse(bIDStr)
				c.Set("booth_id", boothID)
			}
		} else {
			// Kalau user admin, kita simpen user_id-nya
			if uIDStr, ok := claims["user_id"].(string); ok {
				userID, _ := uuid.Parse(uIDStr)
				c.Set("user_id", userID)
			}
		}

		c.Next()
	}
}

// OnlyDevice pastiin cuma mesin fisik (booth) yang bisa tembus.
func OnlyDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")

		if role != "device" {
			response.Error(c, http.StatusForbidden, "Akses ditolak", "Cuma mesin photobooth yang boleh akses")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GenerateDeviceToken dipanggil pas proses pairing buat bikin "Kunci" mesin.
func GenerateDeviceToken(boothID, tenantID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"booth_id":  boothID.String(),
		"tenant_id": tenantID.String(),
		"role":      "device",
		"exp":       time.Now().Add(time.Hour * 24 * 365).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
