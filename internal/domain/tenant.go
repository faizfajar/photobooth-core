package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tenant adalah model data untuk pemilik bisnis (SaaS Owner).
type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantSubscription struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string    `gorm:"type:not null" json:"name"`
	SubscriptionPlan string    `gorm:"type:not null" json:"subscription_plan"`
	Status           string    `gorm:"type:not null" json:"status"`
}

// TenantRepository mendefinisikan cara data disimpan (Database abstraction).
type TenantRepository interface {
	Create(tenant *Tenant) error
}

type TenantSubscriptionRepository interface {
	SubscribePlan(TenantSubscription *TenantSubscription) error
	ChangeSubscribePlan(TenantSubscription *TenantSubscription) error
	UnsubscribePlan(TenantSubscription *TenantSubscription) error
	ChangeStatusPlan(TenantSubscription *TenantSubscription) error
}

// TenantUsecase mendefinisikan aturan bisnis (Business logic abstraction).
type TenantUsecase interface {
	// RegisterTenant(name string) (*Tenant, error)
	RegisterTenant(name, email, password string) (*Tenant, *User, error)
}

type TenantPayment interface {
	PaymentTenant(name string) (*Tenant, error)
}
