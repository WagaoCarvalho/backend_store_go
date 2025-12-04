package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	validate "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func (s *saleService) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if strings.TrimSpace(status) == "" {
		return nil, errMsg.ErrInvalidData
	}

	if err := validate.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validate.ValidateOrder(orderBy, map[string]bool{
		"id":        true,
		"sale_date": true,
		"total":     true,
		"status":    true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByStatus(ctx, status, limit, offset, orderBy, orderDir)
}

// Cancelar venda
func (s *saleService) Cancel(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	saleModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if saleModel.Status != "active" {
		return fmt.Errorf("%w: somente vendas ativas podem ser canceladas", errMsg.ErrInvalidData)
	}

	saleModel.Status = "canceled"
	if err := s.repo.Update(ctx, saleModel); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

// Concluir venda
func (s *saleService) Complete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	saleModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if saleModel.Status != "active" {
		return fmt.Errorf("%w: somente vendas ativas podem ser concluídas", errMsg.ErrInvalidData)
	}

	saleModel.Status = "completed"
	if err := s.repo.Update(ctx, saleModel); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

// Marcar venda como devolvida
func (s *saleService) Returned(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	saleModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if saleModel.Status != "completed" {
		return fmt.Errorf("%w: somente vendas concluídas podem ser devolvidas", errMsg.ErrInvalidData)
	}

	saleModel.Status = "returned"
	if err := s.repo.Update(ctx, saleModel); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

// Reativar venda (transformar em active novamente)
func (s *saleService) Activate(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	saleModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	// Somente vendas canceladas ou devolvidas podem ser reativadas
	if saleModel.Status != "canceled" && saleModel.Status != "returned" {
		return fmt.Errorf("%w: somente vendas canceladas ou devolvidas podem ser reativadas", errMsg.ErrInvalidData)
	}

	saleModel.Status = "active"
	if err := s.repo.Update(ctx, saleModel); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}
