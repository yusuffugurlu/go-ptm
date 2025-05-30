package dtos

type TransactionRequest struct {
	Amount float64 `json:"amount" validate:"required"`
	Type  string  `json:"type" validate:"required"` // "deposit" "withdraw" "transfer"
}

type ScheduledTransactionRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	Date	 string  `json:"date" validate:"required"`
}