// Package handler menangani permintaan HTTP terkait pengelolaan Tenant.
package handler

import (
	"net/http"

	"photobooth-core/internal/domain"

	"github.com/gin-gonic/gin"
)

// TenantHandler adalah struct yang menampung logic handler untuk Tenant.
type TenantHandler struct {
	tenantUsecase domain.TenantUsecase
}

// NewTenantHandler menginisialisasi TenantHandler dengan dependency yang dibutuhkan.
func NewTenantHandler(u domain.TenantUsecase) *TenantHandler {
	return &TenantHandler{
		tenantUsecase: u,
	}
}

// Register menangani pembuatan akun Tenant baru sekaligus akun User Admin pertama.
// Fitur ini merupakan gerbang utama bagi klien SaaS Anda.
func (h *TenantHandler) Register(c *gin.Context) {
	// Definisi struct input untuk binding JSON
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// 1. Melakukan validasi format input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid. Pastikan email benar dan password minimal 6 karakter.",
		})
		return
	}

	// 2. Memanggil logic bisnis di layer Usecase
	// Usecase ini akan menjalankan transaksi database untuk Tenant dan User.
	tenant, user, err := h.tenantUsecase.RegisterTenant(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal melakukan registrasi sistem.",
		})
		return
	}

	// 3. Memberikan respon sukses dengan data Tenant dan User yang terbuat
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi SaaS Berhasil",
		"tenant":  tenant,
		"user":    user,
	})
}
