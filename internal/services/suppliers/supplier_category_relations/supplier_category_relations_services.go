package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
)

var (
	ErrRelationNotFound    = errors.New("relação supplier-categoria não encontrada")
	ErrRelationExists      = errors.New("relação já existe")
	ErrInvalidRelationData = errors.New("dados inválidos para relação")
)

// Interface pública
type SupplierCategoryRelationService interface {
	Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, error)
	GetBySupplier(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategory(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAll(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

// Implementação
type supplierCategoryRelationService struct {
	repository repo.SupplierCategoryRelationRepository
}

// Construtor
func NewSupplierCategoryRelationService(repository repo.SupplierCategoryRelationRepository) SupplierCategoryRelationService {
	return &supplierCategoryRelationService{repository: repository}
}

// Criação da relação
func (s *supplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return nil, ErrInvalidRelationData
	}

	exists, err := s.repository.CheckIfExists(ctx, supplierID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência: %w", err)
	}
	if exists {
		return nil, ErrRelationExists
	}

	relation := &models.SupplierCategoryRelations{
		SupplierID: supplierID,
		CategoryID: categoryID,
	}

	created, err := s.repository.Create(ctx, relation)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar relação: %w", err)
	}

	return created, nil
}

// Busca por fornecedor
func (s *supplierCategoryRelationService) GetBySupplier(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	if supplierID <= 0 {
		return nil, ErrInvalidRelationData
	}
	return s.repository.GetBySupplierID(ctx, supplierID)
}

// Busca por categoria
func (s *supplierCategoryRelationService) GetByCategory(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	if categoryID <= 0 {
		return nil, ErrInvalidRelationData
	}
	return s.repository.GetByCategoryID(ctx, categoryID)
}

// Remove relação específica
func (s *supplierCategoryRelationService) Delete(ctx context.Context, supplierID, categoryID int64) error {
	if supplierID <= 0 || categoryID <= 0 {
		return ErrInvalidRelationData
	}
	if err := s.repository.Delete(ctx, supplierID, categoryID); err != nil {
		if errors.Is(err, repo.ErrRelationNotFound) {
			return ErrRelationNotFound
		}
		return fmt.Errorf("erro ao deletar relação: %w", err)
	}
	return nil
}

// Remove todas as relações de um fornecedor
func (s *supplierCategoryRelationService) DeleteAll(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return ErrInvalidRelationData
	}
	return s.repository.DeleteAllBySupplier(ctx, supplierID)
}

// Verifica se existe uma relação entre supplier e category
func (s *supplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return false, ErrInvalidRelationData
	}
	return s.repository.CheckIfExists(ctx, supplierID, categoryID)
}
