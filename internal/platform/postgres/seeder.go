package postgres

import (
	"log/slog"
	"photobooth-core/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdmin berfungsi untuk membuat akun tenant dan admin pertama jika database masih kosong.
func SeedAdmin(db *gorm.DB) {
	var count int64
	// Cek apakah sudah ada user di database
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		return // Jika sudah ada data, hentikan seeder
	}

	tenantID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)

	tenant := domain.Tenant{
		ID:   tenantID,
		Name: "Admin Utama",
	}

	user := domain.User{
		ID:       uuid.New(),
		TenantID: tenantID,
		Email:    "admin@photobooth.com",
		Password: string(hashedPassword),
	}

	// Menjalankan seeder dalam transaksi database agar aman
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tenant).Error; err != nil {
			return err
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		slog.Error("DATABASE_SEED_FAILED", "error", err)
	} else {
		slog.Info("DATABASE_SEED_SUCCESS", "email", "admin@photobooth.com", "password", "password123")
	}
}
