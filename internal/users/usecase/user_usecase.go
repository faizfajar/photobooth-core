// Package usecase mengimplementasikan logika bisnis untuk autentikasi dan manajemen user.
package usecase

import (
	"errors"
	"photobooth-core/internal/domain"
	"photobooth-core/internal/platform/auth"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase membuat instance baru untuk logika bisnis user.
func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

// Login memverifikasi kredensial dan mengembalikan token JWT jika berhasil.
func (u *userUsecase) Login(email, password string) (string, error) {
	// 1. Mencari user berdasarkan email melalui repository.
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		// Untuk keamanan, jangan spesifikasikan apakah email atau password yang salah.
		return "", errors.New("kredensial yang Anda masukkan salah")
	}

	// 2. Membandingkan password input (plain) dengan password di database (hash).
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("kredensial yang Anda masukkan salah")
	}

	// 3. Membuat token JWT yang mengandung TenantID dan UserID.
	// TenantID sangat penting untuk memfilter data booth secara remote nantinya.
	token, err := auth.GenerateToken(user.TenantID, user.ID)
	if err != nil {
		return "", errors.New("gagal membuat sesi login")
	}

	return token, nil
}
