package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
)

type SupplierCategoryRelationService interface {
	Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, bool, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	DeleteByID(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierID(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelationService struct {
	relationRepo repo.SupplierCategoryRelationRepository
}

func NewSupplierCategoryRelationService(repository repo.SupplierCategoryRelationRepository) SupplierCategoryRelationService {
	return &supplierCategoryRelationService{relationRepo: repository}
}

func (s *supplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, bool, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return nil, false, err_msg.ErrInvalidData
	}

	exists, err := s.relationRepo.HasRelation(ctx, supplierID, categoryID)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	if exists {
		return nil, false, err_msg.ErrRelationExists
	}

	relation := &models.SupplierCategoryRelations{
		SupplierID: supplierID,
		CategoryID: categoryID,
	}

	created, err := s.relationRepo.Create(ctx, relation)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	return created, true, nil
}

func (s *supplierCategoryRelationService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	if supplierID <= 0 {
		return nil, err_msg.ErrInvalidData
	}

	result, err := s.relationRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return result, nil
}

func (s *supplierCategoryRelationService) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	if categoryID <= 0 {
		return nil, err_msg.ErrInvalidData
	}

	result, err := s.relationRepo.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return result, nil
}

func (s *supplierCategoryRelationService) DeleteByID(ctx context.Context, supplierID, categoryID int64) error {
	if supplierID <= 0 || categoryID <= 0 {
		return err_msg.ErrInvalidData
	}

	err := s.relationRepo.Delete(ctx, supplierID, categoryID)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return err_msg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}

func (s *supplierCategoryRelationService) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return err_msg.ErrInvalidData
	}

	err := s.relationRepo.DeleteAllBySupplierID(ctx, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}

func (s *supplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return false, err_msg.ErrInvalidData
	}

	exists, err := s.relationRepo.HasRelation(ctx, supplierID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	return exists, nil
}
