package domain

import "time"

type Brand struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name" validate:"required,min=3,max=200"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type BrandRepository interface {
    Create(brand *Brand) error
    GetByID(id int64) (*Brand, error)
    Update(brand *Brand) error
    Delete(id int64) error
    List() ([]Brand, error)
}

type BrandService interface {
    Create(brand *Brand) error
    GetByID(id int64) (*Brand, error)
    Update(brand *Brand) error
    Delete(id int64) error
    List() ([]Brand, error)
} 