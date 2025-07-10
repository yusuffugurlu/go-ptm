package dtos

type TransactionRequest struct {
	Amount float64 `json:"amount" validate:"required"`
	Type   string  `json:"type" validate:"required"` // "deposit" "withdraw" "transfer"
}

type DebitRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type TransferRequest struct {
	ToUserID uint    `json:"to_user_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required,gt=0"`
}

type TransactionResponse struct {
	ID         uint          `json:"id"`
	FromUserID *uint         `json:"from_user_id,omitempty"`
	ToUserID   *uint         `json:"to_user_id,omitempty"`
	Amount     float64       `json:"amount"`
	Type       string        `json:"type"`
	Status     string        `json:"status"`
	CreatedAt  string        `json:"created_at"`
	FromUser   *UserResponse `json:"from_user,omitempty"`
	ToUser     *UserResponse `json:"to_user,omitempty"`
}

type UserResponse struct {
	ID       uint             `json:"id"`
	Username string           `json:"username"`
	Email    string           `json:"email"`
	Balance  *BalanceResponse `json:"balance,omitempty"`
}

type BalanceResponse struct {
	Amount float64 `json:"amount"`
}

type ScheduledTransactionRequest struct {
	Amount float64 `json:"amount" validate:"required"`
	Date   string  `json:"date" validate:"required"`
}

type HistoricalBalanceResponse struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
