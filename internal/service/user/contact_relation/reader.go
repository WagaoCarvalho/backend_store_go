package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userContactRelationService) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelation, error) {
	if userID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *userContactRelationService) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
	if userID <= 0 {
		return false, errMsg.ErrZeroID
	}
	if contactID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.relationRepo.HasUserContactRelation(ctx, userID, contactID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return exists, nil
}
