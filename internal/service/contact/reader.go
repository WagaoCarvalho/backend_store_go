package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	contact, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %w", errMsg.ErrGet, err)
	}

	return contact, nil
}
