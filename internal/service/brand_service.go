package service

import (
	"api-otto/internal/domain"
)

type brandService struct {
	repository domain.BrandRepository
}

func NewBrandService(repository domain.BrandRepository) domain.BrandService {
	return &brandService{
		repository: repository,
	}
}

func (s *brandService) Create(brand *domain.Brand) error {
	return s.repository.Create(brand)
}

func (s *brandService) GetByID(id int64) (*domain.Brand, error) {
	return s.repository.GetByID(id)
}

func (s *brandService) Update(brand *domain.Brand) error {
	return s.repository.Update(brand)
}

func (s *brandService) Delete(id int64) error {
	return s.repository.Delete(id)
}

func (s *brandService) List() ([]domain.Brand, error) {
	return s.repository.List()
} 