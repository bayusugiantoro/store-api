package domain

import "time"

type Transaction struct {
    ID          int64             `json:"id"`
    CustomerID  int64             `json:"customer_id" validate:"required"`
    TotalPoints int               `json:"total_points"`
    Status      TransactionStatus `json:"status"`
    Items       []TransactionItem `json:"items" validate:"required,min=1"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

type TransactionItem struct {
    ID          int64     `json:"id"`
    TransactionID int64   `json:"transaction_id"`
    VoucherID   int64     `json:"voucher_id" validate:"required"`
    PointsUsed  int       `json:"points_used"`
    CreatedAt   time.Time `json:"created_at"`
    Voucher     *Voucher  `json:"voucher,omitempty"`
}

type TransactionStatus string

const (
    TransactionStatusPending   TransactionStatus = "pending"
    TransactionStatusCompleted TransactionStatus = "completed"
    TransactionStatusFailed    TransactionStatus = "failed"
)

type TransactionRepository interface {
    Create(transaction *Transaction) error
    GetByID(id int64) (*Transaction, error)
    GetByCustomerID(customerID int64) ([]Transaction, error)
    Update(transaction *Transaction) error
    CreateTransactionItem(item *TransactionItem) error
    GetTransactionItems(transactionID int64) ([]TransactionItem, error)
}

type TransactionService interface {
    CreateRedemption(transaction *Transaction) error
    GetTransactionByID(id int64) (*Transaction, error)
    GetCustomerTransactions(customerID int64) ([]Transaction, error)
} 