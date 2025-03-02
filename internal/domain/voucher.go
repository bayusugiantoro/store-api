package domain

import "time"

type Voucher struct {
    ID          int64     `json:"id"`
    BrandID     int64     `json:"brand_id" validate:"required"`
    Code        string    `json:"code" validate:"required"`
    Name        string    `json:"name" validate:"required"`
    Description string    `json:"description"`
    Points      int       `json:"points" validate:"required,gt=0"`
    ValidUntil  time.Time `json:"valid_until"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Brand       *Brand    `json:"brand,omitempty"`
}

type VoucherRepository interface {
    Create(voucher *Voucher) error
    GetByID(id int64) (*Voucher, error)
    GetByBrandID(brandID int64) ([]Voucher, error)
    Update(voucher *Voucher) error
    Delete(id int64) error
    List() ([]Voucher, error)
}

type VoucherService interface {
    Create(voucher *Voucher) error
    GetByID(id int64) (*Voucher, error)
    GetByBrandID(brandID int64) ([]Voucher, error)
    Update(voucher *Voucher) error
    Delete(id int64) error
    List() ([]Voucher, error)
} 