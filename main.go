package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TransactionType işlem türünü belirtir (yatırma veya çekme)
type TransactionType string

const (
	DepositTransaction  TransactionType = "deposit"
	WithdrawTransaction TransactionType = "withdraw"
)

// Transaction bir para yatırma veya çekme işlemini temsil eder
type Transaction struct {
	UserID string          `json:"user_id"`
	Amount float64         `json:"amount"`
	Type   TransactionType `json:"type"`
}

// AccountStore hesap bakiyelerini ve senkronizasyon için bir mutex'i yönetir
type AccountStore struct {
	mu       sync.Mutex
	balances map[string]float64 // Kullanıcı ID'sine göre bakiye
}

// NewAccountStore yeni bir AccountStore oluşturur
func NewAccountStore() *AccountStore {
	return &AccountStore{
		balances: make(map[string]float64),
	}
}

// Deposit bir kullanıcının hesabına para yatırır
func (as *AccountStore) Deposit(userID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("yatırılacak miktar pozitif olmalıdır")
	}

	as.mu.Lock()
	defer as.mu.Unlock()

	as.balances[userID] += amount
	log.Printf("DEPOSIT: UserID: %s, Amount: %.2f, New Balance: %.2f\n", userID, amount, as.balances[userID])
	return nil
}

// Withdraw bir kullanıcının hesabından para çeker
func (as *AccountStore) Withdraw(userID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("çekilecek miktar pozitif olmalıdır")
	}

	as.mu.Lock()
	defer as.mu.Unlock()

	currentBalance, ok := as.balances[userID]
	if !ok {
		// Veya yeni kullanıcı için 0 bakiye varsayılabilir, bu durumda bu kontrol kaldırılır.
		return fmt.Errorf("kullanıcı bulunamadı: %s", userID)
	}

	if currentBalance < amount {
		return fmt.Errorf("yetersiz bakiye: Kullanıcı %s, Bakiye %.2f, Çekilmek İstenen %.2f", userID, currentBalance, amount)
	}

	as.balances[userID] -= amount
	log.Printf("WITHDRAW: UserID: %s, Amount: %.2f, New Balance: %.2f\n", userID, amount, as.balances[userID])
	return nil
}

// GetBalance bir kullanıcının mevcut bakiyesini döndürür (mutex korumalı)
func (as *AccountStore) GetBalance(userID string) (float64, bool) {
	as.mu.Lock()
	defer as.mu.Unlock()
	balance, ok := as.balances[userID]
	return balance, ok
}

// --- Worker Pool ---

const (
	maxWorkers   = 5    // Aynı anda çalışacak maksimum worker sayısı
	maxQueueSize = 100  // İş kuyruğunun maksimum boyutu
)

// JobQueue işlenecek işlemleri tutan kanal
var JobQueue chan Transaction

// WorkerPool worker'ları yönetir
type WorkerPool struct {
	accountStore *AccountStore
}

// NewWorkerPool yeni bir WorkerPool oluşturur ve worker'ları başlatır
func NewWorkerPool(as *AccountStore, numWorkers int) *WorkerPool {
	JobQueue = make(chan Transaction, maxQueueSize)
	wp := &WorkerPool{accountStore: as}

	for i := 1; i <= numWorkers; i++ {
		go wp.worker(i, JobQueue)
	}
	log.Printf("%d adet worker başlatıldı.\n", numWorkers)
	return wp
}

// worker bir goroutine'dir ve JobQueue'dan gelen işlemleri işler
func (wp *WorkerPool) worker(id int, jobs <-chan Transaction) {
	log.Printf("Worker %d başlatıldı ve iş bekliyor...\n", id)
	for job := range jobs {
		log.Printf("Worker %d, UserID %s için %s işlemini aldı (Miktar: %.2f)\n", id, job.UserID, job.Type, job.Amount)
		var err error
		switch job.Type {
		case DepositTransaction:
			err = wp.accountStore.Deposit(job.UserID, job.Amount)
		case WithdrawTransaction:
			err = wp.accountStore.Withdraw(job.UserID, job.Amount)
		default:
			err = fmt.Errorf("bilinmeyen işlem tipi: %s", job.Type)
		}

		if err != nil {
			log.Printf("Worker %d HATA: UserID %s, İşlem %s, Miktar %.2f - Hata: %v\n", id, job.UserID, job.Type, job.Amount, err)
		} else {
			log.Printf("Worker %d BAŞARILI: UserID %s, İşlem %s, Miktar %.2f\n", id, job.UserID, job.Type, job.Amount)
		}
		// Simülasyon için biraz bekleme
		time.Sleep(10000 * time.Millisecond)
	}
	log.Printf("Worker %d durdu.\n", id)
}

// SubmitJob bir işlemi iş kuyruğuna ekler
func SubmitJob(tx Transaction) error {
	select {
	case JobQueue <- tx:
		log.Printf("İşlem kuyruğa eklendi: UserID %s, Tip %s, Miktar %.2f\n", tx.UserID, tx.Type, tx.Amount)
		return nil
	default:
		// Kuyruk doluysa, bu durum ele alınabilir.
		// Örneğin, bir hata döndürülebilir veya işlem bekletilebilir.
		// Bu örnekte basitçe hata döndürüyoruz.
		return fmt.Errorf("işlem kuyruğu dolu, işlem reddedildi: UserID %s", tx.UserID)
	}
}

// --- Echo Handlers ---

type TransactionRequest struct {
	UserID string  `json:"user_id" form:"user_id" query:"user_id"`
	Amount float64 `json:"amount" form:"amount" query:"amount"`
}

func handleDeposit(c echo.Context) error {
	req := new(TransactionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Geçersiz istek: %v", err))
	}

	if req.UserID == "" || req.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "UserID ve pozitif bir Amount gereklidir.")
	}

	tx := Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Type:   DepositTransaction,
	}

	if err := SubmitJob(tx); err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Sprintf("İşlem kuyruğa eklenemedi: %v", err))
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Para yatırma işlemi kuyruğa alındı.",
		"user_id": req.UserID,
		"amount":  fmt.Sprintf("%.2f", req.Amount),
	})
}

func handleWithdraw(c echo.Context) error {
	req := new(TransactionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Geçersiz istek: %v", err))
	}

	if req.UserID == "" || req.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "UserID ve pozitif bir Amount gereklidir.")
	}

	tx := Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Type:   WithdrawTransaction,
	}

	if err := SubmitJob(tx); err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Sprintf("İşlem kuyruğa eklenemedi: %v", err))
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Para çekme işlemi kuyruğa alındı.",
		"user_id": req.UserID,
		"amount":  fmt.Sprintf("%.2f", req.Amount),
	})
}

func handleGetBalance(c echo.Context, as *AccountStore) error {
	userID := c.Param("userID")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "UserID gereklidir.")
	}

	balance, ok := as.GetBalance(userID)
	if !ok {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id": userID,
			"balance": 0.0,
			"message": "Kullanıcı bulunamadı veya henüz işlemi yok.",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	})
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Hesap deposunu oluştur
	accountStore := NewAccountStore()

	// Worker Pool'u başlat
	NewWorkerPool(accountStore, maxWorkers)

	// Routes
	e.POST("/deposit", handleDeposit)
	e.POST("/withdraw", handleWithdraw)
	e.GET("/balance/:userID", func(c echo.Context) error {
		return handleGetBalance(c, accountStore)
	})

	// Kullanıcıların başlangıç bakiyelerini ayarlayabilirsiniz (test için)
	accountStore.Deposit("user1", 1000)
	accountStore.Deposit("user2", 500)


	// Sunucuyu başlat
	log.Println("Echo sunucusu 1323 portunda başlatılıyor...")
	if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}