package service

import (
	"api-otto/internal/domain"
	"errors"
	"time"
)

type transactionService struct {
    repository     domain.TransactionRepository
    voucherRepo    domain.VoucherRepository
}

func NewTransactionService(
    repository domain.TransactionRepository,
    voucherRepo domain.VoucherRepository,
) domain.TransactionService {
    return &transactionService{
        repository:     repository,
        voucherRepo:    voucherRepo,
    }
}

func (s *transactionService) CreateRedemption(transaction *domain.Transaction) error {
    // Validasi items tidak kosong
    if len(transaction.Items) == 0 {
        return errors.New("transaction must have at least one item")
    }

    // Hitung total points dan validasi voucher
    var totalPoints int
    for i, item := range transaction.Items {
        voucher, err := s.voucherRepo.GetByID(item.VoucherID)
        if err != nil {
            return err
        }
        if voucher == nil {
            return errors.New("voucher not found")
        }

        // Validasi voucher masih berlaku
        if !voucher.ValidUntil.IsZero() && voucher.ValidUntil.Before(time.Now()) {
            return errors.New("voucher has expired")
        }

        // Set points yang digunakan sesuai dengan voucher
        transaction.Items[i].PointsUsed = voucher.Points
        totalPoints += voucher.Points
    }

    // Set total points transaksi
    transaction.TotalPoints = totalPoints
    transaction.Status = domain.TransactionStatusPending

    // Buat transaksi
    err := s.repository.Create(transaction)
    if err != nil {
        return err
    }

    // Update status transaksi menjadi completed
    transaction.Status = domain.TransactionStatusCompleted
    return s.repository.Update(transaction)
}

func (s *transactionService) GetTransactionByID(id int64) (*domain.Transaction, error) {
    transaction, err := s.repository.GetByID(id)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Ambil detail items
    items, err := s.repository.GetTransactionItems(transaction.ID)
    if err != nil {
        return nil, err
    }
    transaction.Items = items

    return transaction, nil
}

func (s *transactionService) GetCustomerTransactions(customerID int64) ([]domain.Transaction, error) {
    transactions, err := s.repository.GetByCustomerID(customerID)
    if err != nil {
        return nil, err
    }

    // Ambil detail items untuk setiap transaksi
    for i, transaction := range transactions {
        items, err := s.repository.GetTransactionItems(transaction.ID)
        if err != nil {
            return nil, err
        }
        transactions[i].Items = items
    }

    return transactions, nil
} 