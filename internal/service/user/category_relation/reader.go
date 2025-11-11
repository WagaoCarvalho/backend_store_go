package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userCategoryRelationService) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error) {
	if userID <= 0 {
		return []*models.UserCategoryRelation{}, errMsg.ErrZeroID
	}

	relations, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	// Garantir que nunca retorne nil
	if relations == nil {
		relations = []*models.UserCategoryRelation{}
	}

	return relations, nil
}

func (s *userCategoryRelationService) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		return false, errMsg.ErrZeroID
	}
	if categoryID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.relationRepo.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return exists, nil
}
