package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
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

func setupSupplierContact() (*mockSupplier.MockSupplierContactRelation, *supplierContactRelationHandler) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	mockService := new(mockSupplier.MockSupplierContactRelation)
	handler := NewSupplierContactRelationHandler(mockService, logAdapter)

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

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(r *model.SupplierContactRelation) bool {
			return r.SupplierID == 1 && r.ContactID == 2
		})).Return(expectedModel, nil).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.ContactSupplierRelationDTO `json:"relation"`
		}{Relation: &relationDTO})

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

	t.Run("error - método não permitido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("error - JSON inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - relation não fornecida", func(t *testing.T) {
		_, handler := setupSupplierContact()

		body, _ := json.Marshal(struct {
			Relation *dto.ContactSupplierRelationDTO `json:"relation"`
		}{Relation: nil})

		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - ID zero", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 0,
			ContactID:  2,
		}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrZeroID).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.ContactSupplierRelationDTO `json:"relation"`
		}{Relation: &relationDTO})

		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error - relação já existente", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  2,
		}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrRelationExists).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.ContactSupplierRelationDTO `json:"relation"`
		}{Relation: &relationDTO})

		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error - foreign key inválida", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  99,
		}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.ContactSupplierRelationDTO `json:"relation"`
		}{Relation: &relationDTO})

		req := httptest.NewRequest(http.MethodPost, "/supplier-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error - erro interno", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		relationDTO := dto.ContactSupplierRelationDTO{
			SupplierID: 1,
			ContactID:  2,
		}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error")).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.ContactSupplierRelationDTO `json:"relation"`
		}{Relation: &relationDTO})

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
