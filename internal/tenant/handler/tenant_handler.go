// Package handler menangani permintaan HTTP terkait pengelolaan Tenant.
package handler

import (
	"net/http"

	"photobooth-core/internal/domain"
	"photobooth-core/internal/platform/response"

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

	// Melakukan validasi format input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Validation(c, err)
		return
	}

	// Memanggil logic bisnis di layer Usecase
	// Usecase ini akan menjalankan transaksi database untuk Tenant dan User.
	tenant, user, err := h.tenantUsecase.RegisterTenant(input.Name, input.Email, input.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal registrasi tenant", err.Error())
		return
	}

	// Memberikan respon sukses dengan data Tenant dan User yang terbuat
	response.Success(c, http.StatusCreated, "Registrasi tenant berhasil", gin.H{
		"tenant": tenant,
		"user":   user,
	})
}
