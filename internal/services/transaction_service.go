package services

import (
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/models"
)

type TransactionService interface {
	CreateTransaction(req dtos.TransactionRequest) (*models.Transaction, error)
}

type transactionService struct{}

func NewTransactionService() TransactionService {
	return &transactionService{}
}

func (t *transactionService) CreateTransaction(req dtos.TransactionRequest) (*models.Transaction, error) {
	return nil, nil
}
