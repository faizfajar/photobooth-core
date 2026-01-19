package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(dsn string) (*gorm.DB, error) {
	// Membuka koneksi dengan driver postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Anda bisa menambahkan konfigurasi Connection Pool di sini jika sudah production
	return db, nil
}
