package services

import (
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/repositories"
)

type BalanceService interface {
	GetUserBalance(userID uint) (*dtos.BalanceResponse, error)
}

type balanceService struct {
	balanceRepo repositories.BalancesRepository
}

func NewBalanceService() BalanceService {
	return &balanceService{
		balanceRepo: repositories.NewBalancesRepository(database.Db),
	}
}

func (b *balanceService) GetUserBalance(userID uint) (*dtos.BalanceResponse, error) {
	balance, err := b.balanceRepo.GetByUserId(userID)
	if err != nil {
		return nil, err
	}

	response := &dtos.BalanceResponse{
		Amount: balance.Amount,
	}

	return response, nil
}
