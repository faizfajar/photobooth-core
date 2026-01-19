package domain

import (
	"github.com/google/uuid"
)

// User merepresentasikan tabel 'users'
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID uuid.UUID `gorm:"type:uuid"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"` // Hashed password
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@photobooth.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}
type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"Faiz Photobooth"`
	Email    string `json:"email" binding:"required,email" example:"owner@photobooth.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// UserRepository: Kabel untuk ke Database
type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}

// UserUsecase: Kabel untuk Logika Bisnis
type UserUsecase interface {
	Login(email, password string) (string, error) // Mengembalikan JWT Token
}
