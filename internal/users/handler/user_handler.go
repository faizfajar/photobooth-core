// Package handler menangani permintaan HTTP terkait autentikasi user.
package handler

import (
	"net/http"
	"photobooth-core/internal/domain"
	"photobooth-core/internal/platform/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserHandler membuat instance baru untuk handler user.
func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: u,
	}
}

// Login godoc
// @Summary      User Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      domain.LoginRequest  true  "Kredensial Login"
// @Success      200      {object}  response.Response
// @Failure      401      {object}  response.ErrorResponse
// @Router       /api/v1/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	// validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Validation(c, err)
		return
	}

	// Memanggil logic bisnis di layer Usecase
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login gagal", "Email atau password salah")
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
	})
}
