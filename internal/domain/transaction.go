package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BoothID       uuid.UUID `gorm:"type:uuid;index;not null" json:"booth_id"`
	TenantID      uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ReferenceNo   string    `gorm:"type:varchar(100);unique;not null" json:"reference_no"`
	Amount        float64   `gorm:"type:decimal(10,2)" json:"amount"`
	PaymentStatus string    `gorm:"type:varchar(20);default:'pending'" json:"payment_status"`
	TotalPhotos   int       `gorm:"type:integer;default:0" json:"total_photos"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Booth Booth `gorm:"foreignKey:BoothID" json:"-"`
}

type StartSessionRequest struct {
	ReferenceNo string  `json:"reference_no" binding:"required"`
	Amount      float64 `json:"amount"`
}
