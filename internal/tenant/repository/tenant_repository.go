// Package repository mengimplementasikan akses data untuk entitas Tenant menggunakan PostgreSQL.
package repository

import (
	"photobooth-core/internal/domain"

	"gorm.io/gorm"
)

// tenantRepository adalah implementasi private dari domain.TenantRepository.
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository membuat instance baru untuk manajemen data Tenant.
func NewTenantRepository(db *gorm.DB) domain.TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

// Create menyimpan data Tenant baru ke dalam tabel 'tenants'.
func (r *tenantRepository) Create(tenant *domain.Tenant) error {
	return r.db.Create(tenant).Error
}
