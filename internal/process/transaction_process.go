package process

import (
	"fmt"
	"time"

	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
)

type TransactionType string

const (
	DepositTransaction  TransactionType = "deposit"
	WithdrawTransaction TransactionType = "withdraw"
	TransferTransaction TransactionType = "transfer"
)

type Transaction struct {
	Amount   float32         `json:"amount" validate:"required"`
	UserId   uint            `json:"user_id" validate:"required"`
	ToUserId uint            `json:"to_user_id"`
	Date     time.Time          `json:"date"`
	Type     TransactionType `json:"type" validate:"required"`
}

const (
	maxQueueSize = 100
)

var JobQueue chan Transaction

type WorkerPool struct {
	repo repositories.BalancesRepository
}

func InitWorkerPool(numWorkers int) *WorkerPool {
	JobQueue = make(chan Transaction, maxQueueSize)
	wp := &WorkerPool{
		repo: repositories.NewBalancesRepository(database.Db),
	}

	for i := 1; i <= numWorkers; i++ {
		go wp.worker(i, JobQueue)
	}
	logger.Log.Infof("%d workers started.", numWorkers)
	return wp
}

func (wp *WorkerPool) worker(id int, jobs <-chan Transaction) {
	logger.Log.Infof("Worker %d started and waiting for jobs", id)
	for job := range jobs {
		logger.Log.Infof("Worker %d RECEIVED job for UserID %d: Type %s, Amount %.2f", id, job.UserId, job.Type, job.Amount)
		var err error
		switch job.Type {
		case DepositTransaction:
			wp.repo.Deposit(job.UserId, float64(job.Amount))
		case WithdrawTransaction:
			wp.repo.Withdraw(job.UserId, float64(job.Amount))
		case TransferTransaction:
			err = wp.repo.Transfer(job.UserId, job.ToUserId, float64(job.Amount))
		default:
			err = fmt.Errorf("unknown transaction type: %s", job.Type)
		}

		if err != nil {
			logger.Log.Infof("Worker %d ERROR: UserID %d, Type %s, Amount %.2f - Error: %v", id, job.UserId, job.Type, job.Amount, err)
		} else {
			logger.Log.Infof("Worker %d SUCCESS: UserID %d, Type %s, Amount %.2f", id, job.UserId, job.Type, job.Amount)
		}

		time.Sleep(1000 * time.Millisecond)
	}
	logger.Log.Infof("Worker %d stopped.\n", id)
}

func (wp *WorkerPool) SubmitJob(tx Transaction) error {
	select {
	case JobQueue <- tx:
		logger.Log.Infof("Job added to queue: UserID %d, Type %s, Amount %.2f", tx.UserId, tx.Type, tx.Amount)
		return nil
	default:
		return appErrors.NewConflict(nil, "job queue is full, please try again later")
	}
}
