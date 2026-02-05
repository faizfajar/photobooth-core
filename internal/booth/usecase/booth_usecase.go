package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"photobooth-core/internal/booth/repository"
	"photobooth-core/internal/domain"
	"photobooth-core/internal/middleware"

	"github.com/google/uuid"
)

type BoothUsecase interface {
	RegisterBooth(tenantID uuid.UUID, req domain.CreateBoothRequest) (*domain.Booth, error)
	GetMyBooths(tenantID uuid.UUID) ([]domain.Booth, error)
	PairDevice(req domain.BoothPairingRequest) (*domain.BoothPairingResponse, error)
	Heartbeat(boothID uuid.UUID) error
}

type boothUsecase struct {
	repo repository.BoothRepository
}

func NewBoothUsecase(repo repository.BoothRepository) BoothUsecase {
	return &boothUsecase{repo}
}

func (u *boothUsecase) RegisterBooth(tenantID uuid.UUID, req domain.CreateBoothRequest) (*domain.Booth, error) {
	// Generate random Secret Key
	key := make([]byte, 16)
	rand.Read(key)
	secret := hex.EncodeToString(key)

	// Create a short DeviceCode (e.g., PB-A1B2C3)
	deviceCode := fmt.Sprintf("PB-%s", secret[:6])

	booth := &domain.Booth{
		ID:         uuid.New(),
		TenantID:   tenantID,
		Name:       req.Name,
		DeviceCode: deviceCode,
		SecretKey:  secret,
		Status:     domain.BoothActive,
	}

	if err := u.repo.Create(booth); err != nil {
		return nil, err
	}
	return booth, nil
}

func (u *boothUsecase) GetMyBooths(tenantID uuid.UUID) ([]domain.Booth, error) {
	return u.repo.FindByTenant(tenantID)
}

func (u *boothUsecase) PairDevice(req domain.BoothPairingRequest) (*domain.BoothPairingResponse, error) {
	// 1. Cari booth berdasarkan code yang dikirim mesin
	booth, err := u.repo.FindByDeviceCode(req.DeviceCode)
	if err != nil {
		return nil, fmt.Errorf("device tidak ditemukan")
	}

	// 2. Cek apakah secret key-nya cocok
	if booth.SecretKey != req.SecretKey {
		return nil, fmt.Errorf("secret key salah")
	}

	// 3. Panggil fungsi yang kita buat di middleware tadi
	// Kita nggak butuh parameter "device" lagi karena di dalem fungsinya udah otomatis
	token, err := middleware.GenerateDeviceToken(booth.ID, booth.TenantID)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %v", err)
	}

	// 4. Balikin data buat kebutuhan mesin
	return &domain.BoothPairingResponse{
		Token: token,
		Booth: domain.BoothResponse{
			ID:         booth.ID,
			Name:       booth.Name,
			DeviceCode: booth.DeviceCode,
			Status:     booth.Status,
		},
	}, nil
}

func (u *boothUsecase) Heartbeat(boothID uuid.UUID) error {
	return u.repo.UpdateStatus(boothID, "online")
}
