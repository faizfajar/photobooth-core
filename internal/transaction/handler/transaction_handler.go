package handler

import (
	"net/http"
	"photobooth-core/internal/domain"
	"photobooth-core/internal/platform/response"
	"photobooth-core/internal/transaction/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	usecase usecase.TransactionUsecase
}

func NewTransactionHandler(u usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{u}
}

func (h *TransactionHandler) StartSession(c *gin.Context) {
	// 1. Ambil data dari Token (Injected by Middleware)
	bID, _ := c.Get("booth_id")
	tID, _ := c.Get("tenant_id")

	var req domain.StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Validation(c, err)
		return
	}

	// 2. Eksekusi Usecase
	res, err := h.usecase.CreateSession(bID.(uuid.UUID), tID.(uuid.UUID), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal memulai sesi", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Sesi foto berhasil dicatat", res)
}
