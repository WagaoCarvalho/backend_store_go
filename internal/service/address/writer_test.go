package services

import (
	"context"
	"errors"
	"testing"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressService_Create(t *testing.T) {
	userID := int64(1)

	t.Run("falha quando address é nil", func(t *testing.T) {
		service := NewAddressService(new(mockAddress.MockAddress), nil, nil, nil)

		result, err := service.Create(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrNilModel)
	})

	t.Run("falha na validação do endereço", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		address := &models.Address{}

		result, err := service.Create(context.Background(), address)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("erro do repositório ao criar", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		address := &models.Address{
			UserID:       &userID,
			Street:       "Rua A",
			StreetNumber: "10",
			City:         "Cidade",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "12345678",
			IsActive:     true,
		}

		dbErr := errors.New("db error")
		repo.On("Create", mock.Anything, address).
			Return((*models.Address)(nil), dbErr)

		result, err := service.Create(context.Background(), address)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("sucesso", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		address := &models.Address{
			UserID:       &userID,
			Street:       "Rua A",
			StreetNumber: "10",
			City:         "Cidade",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "12345678",
			IsActive:     true,
		}

		repo.On("Create", mock.Anything, address).
			Return(address, nil)

		result, err := service.Create(context.Background(), address)

		assert.NoError(t, err)
		assert.Equal(t, address, result)
	})
}

func TestAddressService_Update(t *testing.T) {
	userID := int64(1)

	t.Run("falha quando address é nil", func(t *testing.T) {
		service := NewAddressService(new(mockAddress.MockAddress), nil, nil, nil)

		err := service.Update(context.Background(), nil)

		assert.ErrorIs(t, err, errMsg.ErrNilModel)
	})

	t.Run("falha por ID inválido", func(t *testing.T) {
		service := NewAddressService(new(mockAddress.MockAddress), nil, nil, nil)

		err := service.Update(context.Background(), &models.Address{ID: 0})

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("falha na validação", func(t *testing.T) {
		service := NewAddressService(new(mockAddress.MockAddress), nil, nil, nil)

		err := service.Update(context.Background(), &models.Address{
			ID:     1,
			UserID: &userID,
		})

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		address := &models.Address{
			ID:           1,
			UserID:       &userID,
			Street:       "Rua",
			StreetNumber: "1",
			City:         "Cidade",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "12345678",
			IsActive:     true,
		}

		dbErr := errors.New("db error")
		repo.On("Update", mock.Anything, address).Return(dbErr)

		err := service.Update(context.Background(), address)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("sucesso", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		address := &models.Address{
			ID:           1,
			UserID:       &userID,
			Street:       "Rua",
			StreetNumber: "1",
			City:         "Cidade",
			State:        "SP",
			Country:      "Brasil",
			PostalCode:   "12345678",
			IsActive:     true,
		}

		repo.On("Update", mock.Anything, address).Return(nil)

		err := service.Update(context.Background(), address)

		assert.NoError(t, err)
	})
}

func TestAddressService_Delete(t *testing.T) {
	t.Run("falha por ID inválido", func(t *testing.T) {
		service := NewAddressService(new(mockAddress.MockAddress), nil, nil, nil)

		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("retorna ErrNotFound", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		repo.On("Delete", mock.Anything, int64(1)).
			Return(errMsg.ErrNotFound)

		err := service.Delete(context.Background(), 1)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
	})

	t.Run("erro genérico encapsulado em ErrDelete", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		dbErr := errors.New("db error")
		repo.On("Delete", mock.Anything, int64(1)).
			Return(dbErr)

		err := service.Delete(context.Background(), 1)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("sucesso", func(t *testing.T) {
		repo := new(mockAddress.MockAddress)
		service := NewAddressService(repo, nil, nil, nil)

		repo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
	})
}
