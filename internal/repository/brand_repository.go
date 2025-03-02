package repository

import (
	"api-otto/internal/domain"
	"database/sql"
	"time"
)

type brandRepository struct {
    db *sql.DB
}

func NewBrandRepository(db *sql.DB) domain.BrandRepository {
    return &brandRepository{db: db}
}

func (r *brandRepository) Create(brand *domain.Brand) error {
    query := `
        INSERT INTO brands (name, description, created_at, updated_at)
        VALUES ($1, $2, $3, $3)
        RETURNING id`

    now := time.Now()
    return r.db.QueryRow(
        query,
        brand.Name,
        brand.Description,
        now,
    ).Scan(&brand.ID)
}

func (r *brandRepository) GetByID(id int64) (*domain.Brand, error) {
    brand := &domain.Brand{}
    query := `
        SELECT id, name, description, created_at, updated_at
        FROM brands
        WHERE id = $1`

    err := r.db.QueryRow(query, id).Scan(
        &brand.ID,
        &brand.Name,
        &brand.Description,
        &brand.CreatedAt,
        &brand.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return brand, err
}

func (r *brandRepository) List() ([]domain.Brand, error) {
    query := `
        SELECT id, name, description, created_at, updated_at
        FROM brands
        ORDER BY id`

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var brands []domain.Brand
    for rows.Next() {
        var brand domain.Brand
        if err := rows.Scan(
            &brand.ID,
            &brand.Name,
            &brand.Description,
            &brand.CreatedAt,
            &brand.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        brands = append(brands, brand)
    }
    return brands, nil
}

func (r *brandRepository) Update(brand *domain.Brand) error {
    query := `
        UPDATE brands
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4`

    result, err := r.db.Exec(
        query,
        brand.Name,
        brand.Description,
        time.Now(),
        brand.ID,
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

func (r *brandRepository) Delete(id int64) error {
    query := `DELETE FROM brands WHERE id = $1`
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