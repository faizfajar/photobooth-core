package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken membuat JWT token baru untuk user yang berhasil login.
func GenerateToken(tenantID uuid.UUID, userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"tenant_id": tenantID.String(),
		"user_id":   userID.String(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateToken memeriksa apakah token yang dikirim valid.
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing tidak valid")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
