package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booth struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID   uuid.UUID   `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name       string      `gorm:"type:varchar(100);not null" json:"name"`
	DeviceCode string      `gorm:"type:varchar(50);unique;index;not null" json:"device_code"`
	SecretKey  string      `gorm:"type:varchar(100);not null" json:"-"` // Hidden from JSON
	Status     BoothStatus `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`

	// Relationships
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
}

type CreateBoothRequest struct {
	Name string `json:"name" binding:"required" example:"Booth Cabang Sudirman"`
}

// BoothPairingRequest is used when the physical machine first connects
type BoothPairingRequest struct {
	DeviceCode string `json:"device_code" binding:"required"`
	SecretKey  string `json:"secret_key" binding:"required"`
}

type BoothResponse struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	DeviceCode string      `json:"device_code"`
	Status     BoothStatus `json:"status"`
}

type BoothPairingResponse struct {
	Token string        `json:"token"`
	Booth BoothResponse `json:"booth"`
}
