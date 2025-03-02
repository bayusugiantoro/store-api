package repository

import (
	"api-otto/internal/domain"
	"database/sql"
	"time"
)

type transactionRepository struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) domain.TransactionRepository {
    return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *domain.Transaction) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `
        INSERT INTO transactions (customer_id, total_points, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $4)
        RETURNING id`

    now := time.Now()
    err = tx.QueryRow(
        query,
        transaction.CustomerID,
        transaction.TotalPoints,
        transaction.Status,
        now,
    ).Scan(&transaction.ID)
    if err != nil {
        return err
    }

    for _, item := range transaction.Items {
        err = r.createTransactionItemTx(tx, transaction.ID, &item)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *transactionRepository) createTransactionItemTx(tx *sql.Tx, transactionID int64, item *domain.TransactionItem) error {
    query := `
        INSERT INTO transaction_items (transaction_id, voucher_id, points_used, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

    return tx.QueryRow(
        query,
        transactionID,
        item.VoucherID,
        item.PointsUsed,
        time.Now(),
    ).Scan(&item.ID)
}

func (r *transactionRepository) GetByID(id int64) (*domain.Transaction, error) {
    transaction := &domain.Transaction{}
    query := `
        SELECT id, customer_id, total_points, status, created_at, updated_at
        FROM transactions
        WHERE id = $1`

    err := r.db.QueryRow(query, id).Scan(
        &transaction.ID,
        &transaction.CustomerID,
        &transaction.TotalPoints,
        &transaction.Status,
        &transaction.CreatedAt,
        &transaction.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    items, err := r.GetTransactionItems(id)
    if err != nil {
        return nil, err
    }
    transaction.Items = items

    return transaction, nil
}

func (r *transactionRepository) GetByCustomerID(customerID int64) ([]domain.Transaction, error) {
    query := `
        SELECT id, customer_id, total_points, status, created_at, updated_at
        FROM transactions
        WHERE customer_id = $1
        ORDER BY created_at DESC`

    rows, err := r.db.Query(query, customerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []domain.Transaction
    for rows.Next() {
        var transaction domain.Transaction
        if err := rows.Scan(
            &transaction.ID,
            &transaction.CustomerID,
            &transaction.TotalPoints,
            &transaction.Status,
            &transaction.CreatedAt,
            &transaction.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }
    return transactions, nil
}

func (r *transactionRepository) Update(transaction *domain.Transaction) error {
    query := `
        UPDATE transactions
        SET status = $1, updated_at = $2
        WHERE id = $3`

    result, err := r.db.Exec(
        query,
        transaction.Status,
        time.Now(),
        transaction.ID,
    )
    if err != nil {
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (r *transactionRepository) CreateTransactionItem(item *domain.TransactionItem) error {
    query := `
        INSERT INTO transaction_items (transaction_id, voucher_id, points_used, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

    return r.db.QueryRow(
        query,
        item.TransactionID,
        item.VoucherID,
        item.PointsUsed,
        time.Now(),
    ).Scan(&item.ID)
}

func (r *transactionRepository) GetTransactionItems(transactionID int64) ([]domain.TransactionItem, error) {
    query := `
        SELECT ti.id, ti.transaction_id, ti.voucher_id, ti.points_used, ti.created_at,
               v.code, v.name, v.points
        FROM transaction_items ti
        LEFT JOIN vouchers v ON ti.voucher_id = v.id
        WHERE ti.transaction_id = $1`

    rows, err := r.db.Query(query, transactionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []domain.TransactionItem
    for rows.Next() {
        var item domain.TransactionItem
        item.Voucher = &domain.Voucher{}
        
        if err := rows.Scan(
            &item.ID,
            &item.TransactionID,
            &item.VoucherID,
            &item.PointsUsed,
            &item.CreatedAt,
            &item.Voucher.Code,
            &item.Voucher.Name,
            &item.Voucher.Points,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, nil
} 