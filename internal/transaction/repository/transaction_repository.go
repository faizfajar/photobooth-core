package repository

import (
	"photobooth-core/internal/domain"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Save(trx *domain.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Save(trx *domain.Transaction) error {
	return r.db.Create(trx).Error
}
