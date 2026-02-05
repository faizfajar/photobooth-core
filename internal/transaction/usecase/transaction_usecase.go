package usecase

import (
	"photobooth-core/internal/domain"
	"photobooth-core/internal/transaction/repository"
	"time"

	"github.com/google/uuid"
)

type TransactionUsecase interface {
	CreateSession(boothID, tenantID uuid.UUID, req domain.StartSessionRequest) (*domain.Transaction, error)
}

type transactionUsecase struct {
	repo repository.TransactionRepository
}

func NewTransactionUsecase(repo repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{repo}
}

func (u *transactionUsecase) CreateSession(boothID, tenantID uuid.UUID, req domain.StartSessionRequest) (*domain.Transaction, error) {
	trx := &domain.Transaction{
		ID:            uuid.New(),
		BoothID:       boothID,
		TenantID:      tenantID,
		ReferenceNo:   req.ReferenceNo,
		Amount:        req.Amount,
		PaymentStatus: "completed",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := u.repo.Save(trx); err != nil {
		return nil, err
	}

	return trx, nil
}
