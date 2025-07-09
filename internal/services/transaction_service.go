package services

import (
	"time"

	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
)

type TransactionService interface {
	CreateTransaction(req dtos.TransactionRequest) (*models.Transaction, error)
	DebitFromUser(userID uint, amount float64) error
	TransferBetweenUsers(fromUserID, toUserID uint, amount float64) error
	GetTransactionHistory(userID uint, limit, offset int) ([]*dtos.TransactionResponse, error)
	GetTransactionByID(id uint) (*dtos.TransactionResponse, error)
	GetAllTransactions(limit, offset int) ([]*dtos.TransactionResponse, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	balanceRepo     repositories.BalancesRepository
}

func NewTransactionService() TransactionService {
	return &transactionService{
		transactionRepo: repositories.NewTransactionRepository(database.Db),
		balanceRepo:     repositories.NewBalancesRepository(database.Db),
	}
}

func (t *transactionService) CreateTransaction(req dtos.TransactionRequest) (*models.Transaction, error) {
	return nil, nil
}

func (t *transactionService) DebitFromUser(userID uint, amount float64) error {
	if amount <= 0 {
		return appErrors.NewBadRequest(nil, "amount must be positive")
	}

	if err := t.balanceRepo.Deposit(userID, amount); err != nil {
		return err
	}

	transaction := &models.Transaction{
		FromUserId: nil,
		ToUserId:   &userID,
		Amount:     amount,
		Type:       "debit",
		Status:     "completed",
		CreatedAt:  time.Now(),
	}

	return t.transactionRepo.Create(transaction)
}

func (t *transactionService) TransferBetweenUsers(fromUserID, toUserID uint, amount float64) error {
	if amount <= 0 {
		return appErrors.NewBadRequest(nil, "amount must be positive")
	}

	if fromUserID == toUserID {
		return appErrors.NewBadRequest(nil, "cannot transfer to same user")
	}

	return t.transactionRepo.Transfer(fromUserID, toUserID, amount)
}

func (t *transactionService) GetTransactionHistory(userID uint, limit, offset int) ([]*dtos.TransactionResponse, error) {
	transactions, err := t.transactionRepo.GetHistoryByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dtos.TransactionResponse
	for _, tx := range transactions {
		txResponse := &dtos.TransactionResponse{
			ID:         tx.Id,
			FromUserID: tx.FromUserId,
			ToUserID:   tx.ToUserId,
			Amount:     tx.Amount,
			Type:       tx.Type,
			Status:     tx.Status,
			CreatedAt:  tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if tx.FromUser != nil {
			fromUserResponse := &dtos.UserResponse{
				ID:       tx.FromUser.Id,
				Username: tx.FromUser.Username,
				Email:    tx.FromUser.Email,
			}
			if tx.FromUser.Balance != nil {
				fromUserResponse.Balance = &dtos.BalanceResponse{
					Amount: tx.FromUser.Balance.Amount,
				}
			}
			txResponse.FromUser = fromUserResponse
		}

		if tx.ToUser != nil {
			toUserResponse := &dtos.UserResponse{
				ID:       tx.ToUser.Id,
				Username: tx.ToUser.Username,
				Email:    tx.ToUser.Email,
			}
			if tx.ToUser.Balance != nil {
				toUserResponse.Balance = &dtos.BalanceResponse{
					Amount: tx.ToUser.Balance.Amount,
				}
			}
			txResponse.ToUser = toUserResponse
		}

		response = append(response, txResponse)
	}

	return response, nil
}

func (t *transactionService) GetTransactionByID(id uint) (*dtos.TransactionResponse, error) {
	transaction, err := t.transactionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &dtos.TransactionResponse{
		ID:         transaction.Id,
		FromUserID: transaction.FromUserId,
		ToUserID:   transaction.ToUserId,
		Amount:     transaction.Amount,
		Type:       transaction.Type,
		Status:     transaction.Status,
		CreatedAt:  transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if transaction.FromUser != nil {
		fromUserResponse := &dtos.UserResponse{
			ID:       transaction.FromUser.Id,
			Username: transaction.FromUser.Username,
			Email:    transaction.FromUser.Email,
		}
		if transaction.FromUser.Balance != nil {
			fromUserResponse.Balance = &dtos.BalanceResponse{
				Amount: transaction.FromUser.Balance.Amount,
			}
		}
		response.FromUser = fromUserResponse
	}

	if transaction.ToUser != nil {
		toUserResponse := &dtos.UserResponse{
			ID:       transaction.ToUser.Id,
			Username: transaction.ToUser.Username,
			Email:    transaction.ToUser.Email,
		}
		if transaction.ToUser.Balance != nil {
			toUserResponse.Balance = &dtos.BalanceResponse{
				Amount: transaction.ToUser.Balance.Amount,
			}
		}
		response.ToUser = toUserResponse
	}

	return response, nil
}

func (t *transactionService) GetAllTransactions(limit, offset int) ([]*dtos.TransactionResponse, error) {
	transactions, err := t.transactionRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dtos.TransactionResponse
	for _, tx := range transactions {
		txResponse := &dtos.TransactionResponse{
			ID:         tx.Id,
			FromUserID: tx.FromUserId,
			ToUserID:   tx.ToUserId,
			Amount:     tx.Amount,
			Type:       tx.Type,
			Status:     tx.Status,
			CreatedAt:  tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if tx.FromUser != nil {
			fromUserResponse := &dtos.UserResponse{
				ID:       tx.FromUser.Id,
				Username: tx.FromUser.Username,
				Email:    tx.FromUser.Email,
			}
			if tx.FromUser.Balance != nil {
				fromUserResponse.Balance = &dtos.BalanceResponse{
					Amount: tx.FromUser.Balance.Amount,
				}
			}
			txResponse.FromUser = fromUserResponse
		}

		if tx.ToUser != nil {
			toUserResponse := &dtos.UserResponse{
				ID:       tx.ToUser.Id,
				Username: tx.ToUser.Username,
				Email:    tx.ToUser.Email,
			}
			if tx.ToUser.Balance != nil {
				toUserResponse.Balance = &dtos.BalanceResponse{
					Amount: tx.ToUser.Balance.Amount,
				}
			}
			txResponse.ToUser = toUserResponse
		}

		response = append(response, txResponse)
	}

	return response, nil
}
