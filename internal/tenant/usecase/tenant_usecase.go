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

func (u *tenantUsecase) RegisterTenant(name, email, password string) (*domain.Tenant, *domain.User, error) {
	var newTenant *domain.Tenant
	var newUser *domain.User

	// Mulai Transaksi Database
	err := u.db.Transaction(func(tx *gorm.DB) error {
		// Buat Objek Tenant
		newTenant = &domain.Tenant{
			ID:   uuid.New(),
			Name: name,
		}
		if err := u.tenantRepo.Create(newTenant); err != nil {
			return err
		}

		// Hash Password untuk User Admin
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("gagal memproses password")
		}

		// Buat Objek User Admin yang terhubung ke TenantID
		newUser = &domain.User{
			ID:       uuid.New(),
			TenantID: newTenant.ID, // Hubungkan ke Tenant yang baru dibuat
			Email:    email,
			Password: string(hashedPassword),
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
