package usecase

import (
	"errors"
	"photobooth-core/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type tenantUsecase struct {
	tenantRepo domain.TenantRepository
	userRepo   domain.UserRepository
	db         *gorm.DB // Butuh instance DB untuk transaksi
}

// NewTenantUsecase sekarang menerima dua repository.
func NewTenantUsecase(tr domain.TenantRepository, ur domain.UserRepository, db *gorm.DB) domain.TenantUsecase {
	return &tenantUsecase{
		tenantRepo: tr,
		userRepo:   ur,
		db:         db,
	}
}

var newTenant *domain.Tenant

func (u *tenantUsecase) RegisterTenant(req domain.RegisterTenantRequest) (*domain.Tenant, *domain.User, error) {
	var newUser *domain.User

	// Mulai Transaksi Database
	err := u.db.Transaction(func(tx *gorm.DB) error {
		// Buat Objek Tenant
		newTenant = &domain.Tenant{
			ID:   uuid.New(),
			Name: req.TenantName,
		}
		if err := u.tenantRepo.Create(newTenant); err != nil {
			return err
		}

		// Hash Password untuk User Admin
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("gagal memproses password")
		}

		newUser := &domain.User{
			ID:       uuid.New(),
			TenantID: newTenant.ID,
			Name:     req.AdminName,
			Email:    req.Email,
			Password: string(hashedPassword),
			Role:     "admin",
		}
		if err := u.userRepo.Create(newUser); err != nil {
			return err
		}

		// commit transaction
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return newTenant, newUser, nil
}
