package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/supplier"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/contact_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSupplierContact() (*mockSupplier.MockSupplierContactRelationService, *SupplierContactRelation) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	mockService := new(mockSupplier.MockSupplierContactRelationService)
	handler := NewSupplierContactRelation(mockService, logAdapter)

	return mockService, handler
}

func TestSupplierContactRelationHandler_Create(t *testing.T) {
	t.Run("success - relação criada", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  2,
		}
		expectedModel := &model.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
			CreatedAt:  time.Now(),
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(expectedModel, true, nil).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação criada com sucesso", response.Message)

		dataBytes, _ := json.Marshal(response.Data)
		var got dto.ContactSupplierRelationDTO
		_ = json.Unmarshal(dataBytes, &got)

		assert.Equal(t, int64(1), got.SupplierID)
		assert.Equal(t, int64(2), got.ContactID)

		mockService.AssertExpectations(t)
	})

	t.Run("already exists - relação já existente", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  2,
		}
		existingModel := &model.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
			CreatedAt:  time.Now(),
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(existingModel, false, nil).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação já existente", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid foreign key", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  2,
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, false, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  2,
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, false, errors.New("db error")).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierContactRelationHandler_GetAllBySupplierID(t *testing.T) {
	t.Run("success - relações retornadas", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		supplierID := int64(1)
		relations := []*model.SupplierContactRelation{
			{SupplierID: supplierID, ContactID: 2, CreatedAt: time.Now()},
		}

		mockService.On("GetAllRelationsBySupplierID", mock.Anything, supplierID).
			Return(relations, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1"})
		rec := httptest.NewRecorder()

		handler.GetAllBySupplierID(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relações encontradas", response.Message)

		dataBytes, _ := json.Marshal(response.Data)
		var got []dto.ContactSupplierRelationDTO
		_ = json.Unmarshal(dataBytes, &got)

		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), got[0].SupplierID)
		assert.Equal(t, int64(2), got[0].ContactID)

		mockService.AssertExpectations(t)
	})

	t.Run("bad request - id inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetAllBySupplierID(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("internal server error - falha ao buscar relações", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		supplierID := int64(1)
		mockService.On("GetAllRelationsBySupplierID", mock.Anything, supplierID).
			Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1"})
		rec := httptest.NewRecorder()

		handler.GetAllBySupplierID(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierContactRelationHandler_Delete(t *testing.T) {
	t.Run("success - relação deletada", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.On("Delete", mock.Anything, int64(1), int64(2)).
			Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/1/2", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1", "contact_id": "2"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação deletada com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("bad request - supplier_id inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/abc/2", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "abc", "contact_id": "2"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("bad request - contact_id inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/1/xyz", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1", "contact_id": "xyz"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("internal server error - falha ao deletar", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.On("Delete", mock.Anything, int64(1), int64(2)).
			Return(errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/1/2", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1", "contact_id": "2"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierContactRelationHandler_DeleteAll(t *testing.T) {
	t.Run("success - relações deletadas", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.On("DeleteAll", mock.Anything, int64(1)).
			Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1"})
		rec := httptest.NewRecorder()

		handler.DeleteAll(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relações deletadas com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("bad request - supplier_id inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "abc"})
		rec := httptest.NewRecorder()

		handler.DeleteAll(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("internal server error - falha ao deletar", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.On("DeleteAll", mock.Anything, int64(1)).
			Return(errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/supplier-contact-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "1"})
		rec := httptest.NewRecorder()

		handler.DeleteAll(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}
