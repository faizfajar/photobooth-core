package repository

import (
	"photobooth-core/internal/domain"

	"gorm.io/gorm"
)

type usersTenantRepository struct {
	db *gorm.DB
}

// NewUserTenantRepository membuat instance baru untuk repository tenant.
func NewUserTenantRepository(db *gorm.DB) domain.UserRepository {
	return &usersTenantRepository{db}

}

func (r *usersTenantRepository) Create(users *domain.User) error {
	return r.db.Create(users).Error
}

func (r *usersTenantRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
