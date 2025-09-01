package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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
	logger       *logger.LogAdapter
}

func NewSupplierCategoryRelationService(repository repo.SupplierCategoryRelationRepository, logger *logger.LogAdapter) SupplierCategoryRelationService {
	return &supplierCategoryRelationService{relationRepo: repository, logger: logger}
}

func (s *supplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models.SupplierCategoryRelations, bool, error) {
	ref := "[supplierCategoryRelationService - Create] - "

	if supplierID <= 0 || categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return nil, false, err_msg.ErrInvalidData
	}

	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	exists, err := s.relationRepo.HasRelation(ctx, supplierID, categoryID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCheckError, nil)
		return nil, false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	if exists {
		s.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return nil, false, err_msg.ErrRelationExists
	}

	relation := &models.SupplierCategoryRelations{
		SupplierID: supplierID,
		CategoryID: categoryID,
	}

	created, err := s.relationRepo.Create(ctx, relation)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	return created, true, nil
}

func (s *supplierCategoryRelationService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	ref := "[supplierCategoryRelationService - GetBySupplierID] - "

	if supplierID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, err_msg.ErrInvalidData
	}

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
	})

	result, err := s.relationRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id":   supplierID,
		"relations_len": len(result),
	})

	return result, nil
}

func (s *supplierCategoryRelationService) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	ref := "[supplierCategoryRelationService - GetByCategoryID] - "

	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return nil, err_msg.ErrInvalidData
	}

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"category_id": categoryID,
	})

	result, err := s.relationRepo.GetByCategoryID(ctx, categoryID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"category_id": categoryID,
		})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"category_id":   categoryID,
		"relations_len": len(result),
	})

	return result, nil
}

func (s *supplierCategoryRelationService) DeleteByID(ctx context.Context, supplierID, categoryID int64) error {
	ref := "[supplierCategoryRelationService - DeleteByID] - "

	if supplierID <= 0 || categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return err_msg.ErrInvalidData
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	err := s.relationRepo.Delete(ctx, supplierID, categoryID)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": supplierID,
				"category_id": categoryID,
			})
			return err_msg.ErrNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	return nil
}

func (s *supplierCategoryRelationService) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	ref := "[supplierCategoryRelationService - DeleteAllBySupplierID] - "

	if supplierID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
		})
		return err_msg.ErrInvalidData
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"supplier_id": supplierID,
	})

	err := s.relationRepo.DeleteAllBySupplierID(ctx, supplierID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"supplier_id": supplierID,
		})
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": supplierID,
	})

	return nil
}

func (s *supplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	ref := "[supplierCategoryRelationService - HasRelation] - "

	if supplierID <= 0 || categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return false, err_msg.ErrInvalidData
	}

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	exists, err := s.relationRepo.HasRelation(ctx, supplierID, categoryID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
		"exists":      exists,
	})

	return exists, nil
}
