package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
)

var (
	ErrRelationNotFound                        = errors.New("relação supplier-categoria não encontrada")
	ErrRelationExists                          = errors.New("relação já existe")
	ErrInvalidRelationData                     = errors.New("dados inválidos para relação")
	ErrCheckRelationExists                     = errors.New("erro ao verificar existência da relação")
	ErrCreateRelation                          = errors.New("erro ao criar relação")
	ErrDeleteRelation                          = errors.New("erro ao deletar relação")
	ErrDeleteAllRelations                      = errors.New("erro ao deletar todas as relações do fornecedor")
	ErrInvalidSupplierCategoryRelationID       = errors.New("ID da relação é inválido")
	ErrInvalidSupplierCategoryRelationData     = errors.New("dados da relação de categoria do fornecedor são inválidos")
	ErrSupplierCategoryRelationVersionRequired = errors.New("versão da relação é obrigatória")
	ErrSupplierCategoryRelationNotFound        = errors.New("relação de categoria do fornecedor não encontrada")
	ErrSupplierCategoryRelationUpdate          = errors.New("erro ao atualizar a relação de categoria do fornecedor")
)

type SupplierCategoryRelationService interface {
	Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, error)
	GetBySupplierId(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategoryId(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	Update(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
	DeleteById(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierId(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelationService struct {
	repository repo.SupplierCategoryRelationRepository
}

func NewSupplierCategoryRelationService(repository repo.SupplierCategoryRelationRepository) SupplierCategoryRelationService {
	return &supplierCategoryRelationService{repository: repository}
}

func (s *supplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return nil, ErrInvalidRelationData
	}

	exists, err := s.repository.HasSupplierCategoryRelation(ctx, supplierID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCheckRelationExists, err)
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
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	return created, nil
}

func (s *supplierCategoryRelationService) GetBySupplierId(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	if supplierID <= 0 {
		return nil, ErrInvalidRelationData
	}
	return s.repository.GetBySupplierID(ctx, supplierID)
}

func (s *supplierCategoryRelationService) GetByCategoryId(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	if categoryID <= 0 {
		return nil, ErrInvalidRelationData
	}
	return s.repository.GetByCategoryID(ctx, categoryID)
}

func (s *supplierCategoryRelationService) Update(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	if relation.ID <= 0 {
		return nil, ErrInvalidSupplierCategoryRelationID
	}

	if relation.SupplierID <= 0 || relation.CategoryID <= 0 {
		return nil, ErrInvalidRelationData
	}

	if relation.Version <= 0 {
		return nil, ErrSupplierCategoryRelationVersionRequired
	}

	updatedRelation, err := s.repository.Update(ctx, relation)
	if err != nil {
		if errors.Is(err, ErrSupplierCategoryRelationNotFound) {
			return nil, fmt.Errorf("%w", ErrSupplierCategoryRelationNotFound)
		}
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryRelationUpdate, err)
	}

	return updatedRelation, nil
}

func (s *supplierCategoryRelationService) DeleteById(ctx context.Context, supplierID, categoryID int64) error {
	if supplierID <= 0 || categoryID <= 0 {
		return ErrInvalidRelationData
	}
	if err := s.repository.Delete(ctx, supplierID, categoryID); err != nil {
		if errors.Is(err, repo.ErrRelationNotFound) {
			return ErrRelationNotFound
		}
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}
	return nil
}

func (s *supplierCategoryRelationService) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return ErrInvalidRelationData
	}
	if err := s.repository.DeleteAllBySupplierId(ctx, supplierID); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAllRelations, err)
	}
	return nil
}

func (s *supplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return false, ErrInvalidRelationData
	}
	return s.repository.HasSupplierCategoryRelation(ctx, supplierID, categoryID)
}
