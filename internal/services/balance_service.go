package services

import (
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/repositories"
)

type BalanceService interface {
	GetUserBalance(userID uint) (*dtos.BalanceResponse, error)
	GetHistoricalBalances(userID uint) ([]dtos.HistoricalBalanceResponse, error)
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

func (b *balanceService) GetHistoricalBalances(userID uint) ([]dtos.HistoricalBalanceResponse, error) {
	historicalBalances, err := b.balanceRepo.GetHistoricalBalances(userID)
	if err != nil {
		return nil, err
	}

	var response []dtos.HistoricalBalanceResponse
	for _, balance := range historicalBalances {
		response = append(response, dtos.HistoricalBalanceResponse{
			Date:   balance.Date.Format("2006-01-02T15:04:05Z07:00"),
			Amount: balance.Amount,
		})
	}

	return response, nil
}
