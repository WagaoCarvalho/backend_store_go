package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model_address "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	model_supplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	model_categories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	model_supplier_full "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
	service "github.com/WagaoCarvalho/backend_store_go/internal/services/supplier/supplier_full_services/supplier_full_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_CreateFull(t *testing.T) {
	mockService := new(service.MockSupplierFullService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewSupplierFullHandler(mockService, logger)
	cnpj := "00.000.000/0000-00"

	t.Run("Sucesso ao criar fornecedor completo", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedSupplier := &model_supplier_full.SupplierFull{
			Supplier: &model_supplier.Supplier{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: &cnpj,
			},
			Address: &model_address.Address{
				Street: "Rua A",
				City:   "Cidade B",
			},
			Contact: &model_contact.Contact{
				Phone: "123456789",
			},
			Categories: []model_categories.SupplierCategory{
				{ID: 1},
			},
		}

		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"name": "Fornecedor A",
				"cnpj": "12.345.678/0001-90",
			},
			"address": map[string]interface{}{
				"street": "Rua A",
				"city":   "Cidade B",
			},
			"contact": map[string]interface{}{
				"phone": "123456789",
			},
			"categories": []map[string]interface{}{
				{"id": 1},
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("CreateFull",
			mock.Anything, // aceita qualquer tipo que implemente context.Context
			mock.MatchedBy(func(s *model_supplier_full.SupplierFull) bool {
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

		cnpj := "00.000.000/0000-00"

		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"name": "Fornecedor Falha",
				"cnpj": cnpj,
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("CreateFull",
			mock.Anything, // <- Corrigido: agora aceita qualquer context.Context
			mock.MatchedBy(func(s *model_supplier_full.SupplierFull) bool {
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
