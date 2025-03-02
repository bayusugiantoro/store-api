package service

import (
	"api-otto/internal/domain"
	"errors"
	"time"
)

type voucherService struct {
    repository domain.VoucherRepository
    brandRepo  domain.BrandRepository
}

func NewVoucherService(repository domain.VoucherRepository, brandRepo domain.BrandRepository) domain.VoucherService {
    return &voucherService{
        repository: repository,
        brandRepo:  brandRepo,
    }
}

func (s *voucherService) Create(voucher *domain.Voucher) error {
    // Validasi brand exists
    brand, err := s.brandRepo.GetByID(voucher.BrandID)
    if err != nil {
        return err
    }
    if brand == nil {
        return errors.New("brand not found")
    }

    // Validasi valid_until harus di masa depan
    if !voucher.ValidUntil.IsZero() && voucher.ValidUntil.Before(time.Now()) {
        return errors.New("valid_until must be in the future")
    }

    return s.repository.Create(voucher)
}

func (s *voucherService) GetByID(id int64) (*domain.Voucher, error) {
    return s.repository.GetByID(id)
}

func (s *voucherService) GetByBrandID(brandID int64) ([]domain.Voucher, error) {
    // Validasi brand exists
    brand, err := s.brandRepo.GetByID(brandID)
    if err != nil {
        return nil, err
    }
    if brand == nil {
        return nil, errors.New("brand not found")
    }

    return s.repository.GetByBrandID(brandID)
}

func (s *voucherService) Update(voucher *domain.Voucher) error {
    // Validasi brand exists
    brand, err := s.brandRepo.GetByID(voucher.BrandID)
    if err != nil {
        return err
    }
    if brand == nil {
        return errors.New("brand not found")
    }

    // Validasi voucher exists
    existing, err := s.repository.GetByID(voucher.ID)
    if err != nil {
        return err
    }
    if existing == nil {
        return errors.New("voucher not found")
    }

    return s.repository.Update(voucher)
}

func (s *voucherService) Delete(id int64) error {
    return s.repository.Delete(id)
}

func (s *voucherService) List() ([]domain.Voucher, error) {
    return s.repository.List()
} 