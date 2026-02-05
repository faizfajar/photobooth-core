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

// Register godoc
// @Summary      Registrasi Tenant & Admin Baru
// @Description  Membuat perusahaan (Tenant) sekaligus akun Admin pertama secara atomic
// @Tags         Tenants
// @Accept       json
// @Produce      json
// @Param        request  body      domain.RegisterTenantRequest  true  "Data Registrasi"
// @Success      201      {object}  response.Response
// @Router       /api/v1/tenants [post]
func (h *TenantHandler) Register(c *gin.Context) {
	//Inisialisasi struct request dari domain
	var req domain.RegisterTenantRequest

	// Jika format JSON salah atau field yang 'required' tidak ada, ini akan error
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Validation(c, err)
		return
	}

	tenant, user, err := h.tenantUsecase.RegisterTenant(req)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal registrasi tenant", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Tenant dan Admin berhasil didaftarkan", gin.H{
		"tenant": tenant,
		"admin":  user,
	})
}
