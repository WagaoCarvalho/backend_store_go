package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_contact_relation"
)

type SupplierContactRelation interface {
	Create(ctx context.Context, supplierID, contactID int64) (*models.SupplierContactRelation, bool, error)
	GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelation, error)
	HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error)
	Delete(ctx context.Context, supplierID, contactID int64) error
	DeleteAll(ctx context.Context, supplierID int64) error
}

type supplierContactRelation struct {
	relationRepo repo.SupplierContactRelation
}

func NewSupplierContactRelation(repo repo.SupplierContactRelation) SupplierContactRelation {
	return &supplierContactRelation{
		relationRepo: repo,
	}
}

func (s *supplierContactRelation) Create(ctx context.Context, supplierID, contactID int64) (*models.SupplierContactRelation, bool, error) {
	if supplierID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}
	if contactID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}

	relation := models.SupplierContactRelation{
		SupplierID: supplierID,
		ContactID:  contactID,
	}

	createdRelation, err := s.relationRepo.Create(ctx, &relation)
	if err != nil {
		switch {
		case errors.Is(err, err_msg.ErrRelationExists):
			relations, getErr := s.relationRepo.GetAllRelationsBySupplierID(ctx, supplierID)
			if getErr != nil {
				return nil, false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, getErr)
			}

			for _, rel := range relations {
				if rel.ContactID == contactID {
					return rel, false, nil
				}
			}

			return nil, false, err_msg.ErrRelationExists

		case errors.Is(err, err_msg.ErrDBInvalidForeignKey):
			return nil, false, err_msg.ErrDBInvalidForeignKey

		default:
			return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
		}
	}

	return createdRelation, true, nil
}

func (s *supplierContactRelation) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelation, error) {
	if supplierID <= 0 {
		return nil, err_msg.ErrZeroID
	}

	relations, err := s.relationRepo.GetAllRelationsBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return relations, nil
}

func (s *supplierContactRelation) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
	if supplierID <= 0 || contactID <= 0 {
		return false, err_msg.ErrZeroID
	}

	exists, err := s.relationRepo.HasSupplierContactRelation(ctx, supplierID, contactID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return exists, nil
}

func (s *supplierContactRelation) Delete(ctx context.Context, supplierID, contactID int64) error {
	if supplierID <= 0 || contactID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, supplierID, contactID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}

func (s *supplierContactRelation) DeleteAll(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
