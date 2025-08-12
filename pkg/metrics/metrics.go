package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UserRegistrationTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_registration_total",
		Help: "Total number of user registrations",
	})

	UserLoginTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_login_total",
		Help: "Total number of user logins",
	})

	TransactionTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "transaction_total",
		Help: "Total number of transactions",
	})

	TransactionAmount = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "transaction_amount",
		Help:    "Transaction amounts distribution",
		Buckets: prometheus.DefBuckets,
	})

	BalanceUpdateTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "balance_update_total",
		Help: "Total number of balance updates",
	})

	DatabaseConnectionStatus = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "database_connection_status",
		Help: "Database connection status (1 = connected, 0 = disconnected)",
	})

	ActiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_users",
		Help: "Number of currently active users",
	})

	ErrorRate = promauto.NewCounter(prometheus.CounterOpts{
		Name: "error_total",
		Help: "Total number of errors",
	})
)

func IncrementUserRegistration() {
	UserRegistrationTotal.Inc()
}

func IncrementUserLogin() {
	UserLoginTotal.Inc()
}

func IncrementTransaction() {
	TransactionTotal.Inc()
}

func ObserveTransactionAmount(amount float64) {
	TransactionAmount.Observe(amount)
}

func IncrementBalanceUpdate() {
	BalanceUpdateTotal.Inc()
}

func SetDatabaseConnectionStatus(status float64) {
	DatabaseConnectionStatus.Set(status)
}

func SetActiveUsers(count float64) {
	ActiveUsers.Set(count)
}

func IncrementError() {
	ErrorRate.Inc()
}
