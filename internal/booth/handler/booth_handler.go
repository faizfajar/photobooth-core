package handler

import (
	"net/http"
	"photobooth-core/internal/booth/usecase"
	"photobooth-core/internal/domain"
	"photobooth-core/internal/platform/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BoothHandler struct {
	usecase usecase.BoothUsecase
}

func NewBoothHandler(u usecase.BoothUsecase) *BoothHandler {
	return &BoothHandler{u}
}

// Register godoc
// @Summary      Register a new Booth
// @Tags         Booths
// @Security     BearerAuth
// @Param        request body domain.CreateBoothRequest true "Booth Data"
// @Success      201 {object} response.Response
// @Router       /api/v1/booths [post]
func (h *BoothHandler) Register(c *gin.Context) {
	var req domain.CreateBoothRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Validation(c, err)
		return
	}

	// Ambil dari context (ini masih interface{})
	tenantIDRaw, _ := c.Get("tenant_id")

	tenantID, err := uuid.Parse(tenantIDRaw.(string)) // Cast ke string dulu baru di-Parse
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Format Tenant ID tidak valid", err.Error())
		return
	}

	// 3. Sekarang 'tenantID' sudah bertipe uuid.UUID
	res, err := h.usecase.RegisterBooth(tenantID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to register booth", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Booth registered successfully", res)
}

// GetAllBooth godoc
// @Summary      Get All Booths for Tenant
// @Tags         Booths
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Router       /api/v1/booths [get]
func (h *BoothHandler) GetAllBooth(c *gin.Context) {
	// get tenant_id by context injected on auth middleware
	tenantIDRaw, exists := c.Get("tenant_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "Tenant ID missing")
		return
	}

	// Type Assertion secara aman
	tenantID, ok := uuid.Parse(tenantIDRaw.(string))
	if ok != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error", "Invalid tenant_id type in context")
		return
	}

	// 3. Panggil Usecase
	booths, err := h.usecase.GetMyBooths(tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data booth", err.Error())
		return
	}

	// 4. Return data (walaupun kosong, kirim [] bukan null)
	response.Success(c, http.StatusOK, "Berhasil mengambil daftar booth", booths)
}

// Pair godoc
// @Summary      Device Handshake (Pairing)
// @Description  Endpoint khusus untuk mesin fisik melakukan login menggunakan Device Code & Secret Key
// @Tags         Booths
// @Accept       json
// @Produce      json
// @Param        request  body      domain.BoothPairingRequest  true  "Pairing Data"
// @Success      200      {object}  response.Response
// @Failure      401      {object}  response.ErrorResponse
// @Router       /api/v1/booths/pair [post]
func (h *BoothHandler) Pair(c *gin.Context) {
	var req domain.BoothPairingRequest

	// validate input json
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Validation(c, err)
		return
	}

	res, err := h.usecase.PairDevice(req)
	if err != nil {
		// Jika device tidak ditemukan atau secret salah, kirim 401 Unauthorized
		response.Error(c, http.StatusUnauthorized, "Pairing gagal", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Device berhasil dipasangkan", res)
}

func (h *BoothHandler) Heartbeat(c *gin.Context) {
	// Ambil ID dari token mesin
	bID, _ := c.Get("booth_id")

	if err := h.usecase.Heartbeat(bID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal update status", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Booth is alive", nil)
}
