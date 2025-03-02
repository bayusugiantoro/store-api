package repository

import (
	"api-otto/internal/domain"
	"database/sql"
	"time"
)

type voucherRepository struct {
    db *sql.DB
}

func NewVoucherRepository(db *sql.DB) domain.VoucherRepository {
    return &voucherRepository{db: db}
}

func (r *voucherRepository) Create(voucher *domain.Voucher) error {
    query := `
        INSERT INTO vouchers (brand_id, code, name, description, points, valid_until, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
        RETURNING id`

    now := time.Now()
    return r.db.QueryRow(
        query,
        voucher.BrandID,
        voucher.Code,
        voucher.Name,
        voucher.Description,
        voucher.Points,
        voucher.ValidUntil,
        now,
    ).Scan(&voucher.ID)
}

func (r *voucherRepository) GetByID(id int64) (*domain.Voucher, error) {
    query := `
        SELECT v.id, v.brand_id, v.code, v.name, v.description, v.points, v.valid_until, 
               v.created_at, v.updated_at,
               b.id, b.name, b.description
        FROM vouchers v
        LEFT JOIN brands b ON v.brand_id = b.id
        WHERE v.id = $1`

    voucher := &domain.Voucher{
        Brand: &domain.Brand{},
    }

    err := r.db.QueryRow(query, id).Scan(
        &voucher.ID,
        &voucher.BrandID,
        &voucher.Code,
        &voucher.Name,
        &voucher.Description,
        &voucher.Points,
        &voucher.ValidUntil,
        &voucher.CreatedAt,
        &voucher.UpdatedAt,
        &voucher.Brand.ID,
        &voucher.Brand.Name,
        &voucher.Brand.Description,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return voucher, err
}

func (r *voucherRepository) GetByBrandID(brandID int64) ([]domain.Voucher, error) {
    query := `
        SELECT v.id, v.brand_id, v.code, v.name, v.description, v.points, v.valid_until, 
               v.created_at, v.updated_at
        FROM vouchers v
        WHERE v.brand_id = $1
        ORDER BY v.id`

    rows, err := r.db.Query(query, brandID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var vouchers []domain.Voucher
    for rows.Next() {
        var voucher domain.Voucher
        if err := rows.Scan(
            &voucher.ID,
            &voucher.BrandID,
            &voucher.Code,
            &voucher.Name,
            &voucher.Description,
            &voucher.Points,
            &voucher.ValidUntil,
            &voucher.CreatedAt,
            &voucher.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        vouchers = append(vouchers, voucher)
    }
    return vouchers, nil
}

func (r *voucherRepository) List() ([]domain.Voucher, error) {
    query := `
        SELECT v.id, v.brand_id, v.code, v.name, v.description, v.points, v.valid_until, 
               v.created_at, v.updated_at
        FROM vouchers v
        ORDER BY v.id`

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var vouchers []domain.Voucher
    for rows.Next() {
        var voucher domain.Voucher
        if err := rows.Scan(
            &voucher.ID,
            &voucher.BrandID,
            &voucher.Code,
            &voucher.Name,
            &voucher.Description,
            &voucher.Points,
            &voucher.ValidUntil,
            &voucher.CreatedAt,
            &voucher.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        vouchers = append(vouchers, voucher)
    }
    return vouchers, nil
}

func (r *voucherRepository) Update(voucher *domain.Voucher) error {
    query := `
        UPDATE vouchers
        SET brand_id = $1, code = $2, name = $3, description = $4, 
            points = $5, valid_until = $6, updated_at = $7
        WHERE id = $8`

    result, err := r.db.Exec(
        query,
        voucher.BrandID,
        voucher.Code,
        voucher.Name,
        voucher.Description,
        voucher.Points,
        voucher.ValidUntil,
        time.Now(),
        voucher.ID,
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

func (r *voucherRepository) Delete(id int64) error {
    query := `DELETE FROM vouchers WHERE id = $1`
    result, err := r.db.Exec(query, id)
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