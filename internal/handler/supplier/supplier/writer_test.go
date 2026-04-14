package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_Create(t *testing.T) {
	t.Run("successfully create supplier", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSupplier := &models.Supplier{
			ID:     1,
			Name:   "Fornecedor Teste",
			CNPJ:   utils.StrToPtr("12345678000199"),
			Status: true,
		}

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor Teste",
			CNPJ: utils.StrToPtr("12345678000199"),
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.Name == "Fornecedor Teste" && s.CNPJ != nil && *s.CNPJ == "12345678000199"
		})).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedor criado com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when JSON is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "JSON inválido", response.Message)

		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return bad request when request body is empty", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader([]byte("")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "JSON inválido", response.Message)

		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return unprocessable entity when validation fails - empty name", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "",
			CNPJ: utils.StrToPtr("12345678000199"),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return unprocessable entity when validation fails - name too long", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		longName := strings.Repeat("a", 101)
		requestBody := dto.SupplierDTO{
			Name: longName,
			CNPJ: utils.StrToPtr("12345678000199"),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return unprocessable entity when CPF and CNPJ both provided", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor Inválido",
			CNPJ: utils.StrToPtr("12345678000199"),
			CPF:  utils.StrToPtr("12345678901"),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return unprocessable entity when CNPJ is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor CNPJ Inválido",
			CNPJ: utils.StrToPtr("12345678"),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return unprocessable entity when CPF is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor CPF Inválido",
			CPF:  utils.StrToPtr("12345678"),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("return internal server error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor Teste",
			CNPJ: utils.StrToPtr("12345678000199"),
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return conflict when supplier already exists (duplicate CNPJ)", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor Duplicado",
			CNPJ: utils.StrToPtr("12345678000199"),
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.Anything).Return(nil, errMsg.ErrDuplicate).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "fornecedor já existente", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return conflict when supplier already exists (duplicate CPF)", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor CPF Duplicado",
			CPF:  utils.StrToPtr("12345678901"),
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.Anything).Return(nil, errMsg.ErrDuplicate).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "fornecedor já existente", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("successfully create supplier with CPF", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSupplier := &models.Supplier{
			ID:     2,
			Name:   "Fornecedor PF",
			CPF:    utils.StrToPtr("12345678901"),
			Status: true,
		}

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor PF",
			CPF:  utils.StrToPtr("12345678901"),
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.Name == "Fornecedor PF" && s.CPF != nil && *s.CPF == "12345678901"
		})).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedor criado com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("successfully create supplier without CPF and CNPJ", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSupplier := &models.Supplier{
			ID:     3,
			Name:   "Fornecedor Sem Documento",
			Status: true,
		}

		requestBody := dto.SupplierDTO{
			Name: "Fornecedor Sem Documento",
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.Name == "Fornecedor Sem Documento" && s.CNPJ == nil && s.CPF == nil
		})).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedor criado com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return unprocessable entity when is_active is provided with invalid value", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		// is_active não é usado no create, mas se fornecido, não deve quebrar
		requestBody := map[string]interface{}{
			"name":      "Fornecedor Teste",
			"cnpj":      "12345678000199",
			"is_active": "invalid", // Tipo inválido
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		// Deve retornar bad request por JSON inválido (tipo mismatch)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertNotCalled(t, "Create")
	})
}
func TestSupplierHandler_Update(t *testing.T) {
	t.Run("successfully update supplier", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Atualizado",
			CNPJ:    utils.StrToPtr("12345678000199"),
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 &&
				s.Name == "Fornecedor Atualizado" &&
				s.Version == 1
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedor atualizado com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when id is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Teste",
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/suppliers/abc", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("return bad request when JSON is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte("{invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("return unprocessable entity when validation fails - empty name", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "",
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("return unprocessable entity when validation fails - name too long", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		longName := strings.Repeat("a", 101)
		requestBody := dto.SupplierDTO{
			Name:    longName,
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("return not found when supplier does not exist", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Não Encontrado",
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 999 && s.Name == "Fornecedor Não Encontrado"
		})).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/999", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "fornecedor não encontrado", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return conflict when version mismatch", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Conflito",
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Version == 1
		})).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "conflito de versão")

		mockService.AssertExpectations(t)
	})

	t.Run("return conflict when supplier duplicate", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Duplicado",
			CNPJ:    utils.StrToPtr("12345678000199"),
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Name == "Fornecedor Duplicado"
		})).Return(errMsg.ErrDuplicate).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "fornecedor já existente", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Erro",
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1
		})).Return(errors.New("database connection error")).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when CPF and CNPJ both provided", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		requestBody := dto.SupplierDTO{
			Name:    "Fornecedor Inválido",
			CNPJ:    utils.StrToPtr("12345678000199"),
			CPF:     utils.StrToPtr("12345678901"),
			Version: 1,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		mockService.AssertNotCalled(t, "Update")
	})
}

func TestSupplierHandler_Delete(t *testing.T) {
	t.Run("successfully delete supplier", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when id is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return not found when supplier does not exist", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("Delete", mock.Anything, int64(999)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
