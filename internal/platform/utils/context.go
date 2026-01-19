package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTenantID mengambil ID tenant dari context yang di-set oleh middleware
func GetTenantID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, errors.New("tenant_id tidak ditemukan di context")
	}

	// Konversi interface ke string, lalu string ke uuid.UUID
	idStr, ok := val.(string)
	if !ok {
		return uuid.Nil, errors.New("format tenant_id tidak valid")
	}

	return uuid.Parse(idStr)
}
