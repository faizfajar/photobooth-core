package middleware

import (
	"net/http"
	"strings"

	"photobooth-core/internal/platform/auth"
	"photobooth-core/internal/platform/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware bertugas memvalidasi token JWT di setiap request terproteksi.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Ambil header Authorization (format: Bearer <token>)
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.ERROR(c, http.StatusUnauthorized, "Otorisasi diperlukan", nil)
			c.Abort()
			return
		}

		//Ekstrak token dari string "Bearer "
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format header otorisasi salah"})
			c.Abort()
			return
		}

		token, err := auth.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			response.ERROR(c, http.StatusUnauthorized, "Sesi berakhir, silakan login kembali", err.Error())
			c.Abort()
			return
		}

		//Perbaikan Type Assertion: Gunakan jwt.MapClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal memproses klaim token"})
			c.Abort()
			return
		}

		//Simpan ke Context (jwt.MapClaims adalah map[string]interface{} di balik layar)
		c.Set("tenant_id", claims["tenant_id"])
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}
