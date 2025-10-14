package services

import (
	"context"
	"time"

	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/process"
	"github.com/yusuffugurlu/go-project/internal/repositories"
)

type ScheduledTransactionService interface {
	Schedule(userId uint, amount float64, at time.Time) (*models.ScheduledTransaction, error)
	StartScheduler(ctx context.Context, interval time.Duration)
}

type scheduledTransactionService struct {
	repo repositories.ScheduledTransactionRepository
}

func NewScheduledTransactionService() ScheduledTransactionService {
	return &scheduledTransactionService{
		repo: repositories.NewScheduledTransactionRepository(database.Db),
	}
}

func (s *scheduledTransactionService) Schedule(userId uint, amount float64, at time.Time) (*models.ScheduledTransaction, error) {
	st := &models.ScheduledTransaction{
		UserId:      userId,
		Amount:      amount,
		ScheduledAt: at,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *scheduledTransactionService) StartScheduler(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				logger.Log.Info("Scheduled transaction scheduler stopped")
				return
			case now := <-ticker.C:
				due, err := s.repo.GetDue(now)
				if err != nil {
					logger.Log.Error("Failed to query due scheduled transactions", err)
					continue
				}
				for _, st := range due {
					process.JobQueue <- process.Transaction{
						Amount: float32(st.Amount),
						UserId: st.UserId,
						Type:   process.DebitTransaction,
						Date:   st.ScheduledAt,
					}
					if err := s.repo.MarkProcessed(st.Id); err != nil {
						logger.Log.Error("Failed to mark scheduled transaction processed", err)
					}
				}
			}
		}
	}()
}
