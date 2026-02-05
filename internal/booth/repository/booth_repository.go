package repository

import (
	"photobooth-core/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoothRepository interface {
	Create(booth *domain.Booth) error
	FindByTenant(tenantID uuid.UUID) ([]domain.Booth, error)
	FindByDeviceCode(code string) (*domain.Booth, error)
	UpdateStatus(id uuid.UUID, status string) error
}

type boothRepository struct {
	db *gorm.DB
}

func NewBoothRepository(db *gorm.DB) BoothRepository {
	return &boothRepository{db}
}

func (r *boothRepository) Create(booth *domain.Booth) error {
	return r.db.Create(booth).Error
}

func (r *boothRepository) FindByTenant(tenantID uuid.UUID) ([]domain.Booth, error) {
	var booths []domain.Booth
	err := r.db.Where("tenant_id = ?", tenantID).Find(&booths).Error
	return booths, err
}

func (r *boothRepository) FindByDeviceCode(code string) (*domain.Booth, error) {
	var booth domain.Booth
	err := r.db.Where("device_code = ?", code).First(&booth).Error
	return &booth, err
}

func (r *boothRepository) UpdateStatus(id uuid.UUID, status string) error {
	// Kita update status dan timestamp 'updated_at' otomatis oleh GORM
	return r.db.Model(&domain.Booth{}).Where("id = ?", id).Update("status", status).Error
}
