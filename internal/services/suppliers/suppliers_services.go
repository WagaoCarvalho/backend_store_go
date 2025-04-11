package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers"
)

type SupplierService interface {
	Create(ctx context.Context, supplier *models.Supplier) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
}

type supplierService struct {
	repo repository.SupplierRepository
}

func NewSupplierService(repo repository.SupplierRepository) SupplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) Create(ctx context.Context, supplier *models.Supplier) (int64, error) {
	if supplier.Name == "" {
		return 0, fmt.Errorf("nome do fornecedor é obrigatório")
	}
	return s.repo.Create(ctx, supplier)
}

func (s *supplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *supplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	return s.repo.GetAll(ctx)
}

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) error {
	if supplier.Name == "" {
		return fmt.Errorf("nome do fornecedor é obrigatório")
	}
	return s.repo.Update(ctx, supplier)
}

func (s *supplierService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return fmt.Errorf("ID inválido para deletar fornecedor")
	}
	return s.repo.Delete(ctx, id)
}
