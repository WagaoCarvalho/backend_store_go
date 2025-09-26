package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"
)

type ClientService interface {
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	GetByName(ctx context.Context, name string) ([]*models.Client, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
	GetAll(ctx context.Context) ([]*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
	ClientExists(ctx context.Context, clientID int64) (bool, error)
}

type clientService struct {
	repo repo.ClientRepository
}

func NewClientService(repo repo.ClientRepository) ClientService {
	return &clientService{
		repo: repo,
	}
}

func (s *clientService) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	created, err := s.repo.Create(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return created, nil
}

func (s *clientService) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return client, nil
}

func (s *clientService) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
	if name == "" {
		return nil, errMsg.ErrInvalidData
	}

	clients, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return clients, nil
}

func (s *clientService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return version, nil
}

func (s *clientService) GetAll(ctx context.Context) ([]*models.Client, error) {
	clients, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return clients, nil
}

func (s *clientService) Update(ctx context.Context, client *models.Client) error {
	if client.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if client.Version <= 0 {
		return errMsg.ErrVersionConflict
	}

	if err := client.Validate(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (s *clientService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}
	return nil
}

func (s *clientService) Disable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Disable(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (s *clientService) Enable(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Enable(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}
	return nil
}

func (s *clientService) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	if clientID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.repo.ClientExists(ctx, clientID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return exists, nil
}
