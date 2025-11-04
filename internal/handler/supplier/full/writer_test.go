package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplierFull "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	dtoSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category"
	dtoSupplierFull "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/full"
	dtoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
	modelSupplier_full "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/full"
	modelSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_CreateFull(t *testing.T) {
	mockService := new(mockSupplierFull.MockSupplierFull)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierFull(mockService, logger)
	cnpj := "00.000.000/0000-00"

	t.Run("Sucesso ao criar fornecedor completo", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedSupplier := &modelSupplier_full.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: &cnpj,
			},
			Address: &modelAddress.Address{
				Street: "Rua A",
				City:   "Cidade B",
			},
			Contact: &modelContact.Contact{
				Phone: "123456789",
			},
			Categories: []modelCategories.SupplierCategory{
				{ID: 1},
			},
		}

		requestDTO := dtoSupplierFull.SupplierFullDTO{
			Supplier: &dtoSupplier.SupplierDTO{
				Name: "Fornecedor A",
				CNPJ: utils.StrToPtr("12.345.678/0001-90"),
			},
			Address: &dtoAddress.AddressDTO{
				Street: "Rua A",
				City:   "Cidade B",
			},
			Contact: &dtoContact.ContactDTO{
				Phone: "123456789",
			},
			Categories: []dtoSupplierCategories.SupplierCategoryDTO{
				{ID: utils.Int64Ptr(1)},
			},
		}

		body, _ := json.Marshal(requestDTO)

		mockService.On("CreateFull",
			mock.Anything,
			mock.MatchedBy(func(s *modelSupplier_full.SupplierFull) bool {
				return s.Supplier.Name == "Fornecedor A" &&
					s.Supplier.CNPJ != nil &&
					*s.Supplier.CNPJ == "12.345.678/0001-90" &&
					s.Address != nil && s.Address.Street == "Rua A" &&
					s.Contact != nil && s.Contact.Phone == "123456789" &&
					len(s.Categories) == 1 && s.Categories[0].ID == 1
			}),
		).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers/full", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers/full", nil)
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ao decodificar JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/suppliers/full", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao criar fornecedor completo no service", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		requestDTO := dtoSupplierFull.SupplierFullDTO{
			Supplier: &dtoSupplier.SupplierDTO{
				Name: "Fornecedor Falha",
				CNPJ: utils.StrToPtr(cnpj),
			},
		}

		body, _ := json.Marshal(requestDTO)

		mockService.On("CreateFull",
			mock.Anything,
			mock.MatchedBy(func(s *modelSupplier_full.SupplierFull) bool {
				return s != nil &&
					s.Supplier != nil &&
					s.Supplier.Name == "Fornecedor Falha" &&
					s.Supplier.CNPJ != nil &&
					*s.Supplier.CNPJ == cnpj
			}),
		).Return(nil, errors.New("erro ao criar fornecedor completo")).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers/full", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
