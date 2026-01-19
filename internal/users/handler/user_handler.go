// Package handler menangani permintaan HTTP terkait autentikasi user.
package handler

import (
	"net/http"
	"photobooth-core/internal/domain"

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
	var req struct {
		Name     string `json:"name" binding:"required,name"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password tidak valid"})
		return
	}

	// Memanggil logic bisnis di layer Usecase
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Mengembalikan token JWT yang berisi TenantID untuk isolasi data
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
