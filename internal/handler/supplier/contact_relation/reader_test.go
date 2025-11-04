package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/contact_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestSupplierContactRelationHandler_HasSupplierContactRelation(t *testing.T) {
	t.Run("success - relação existe", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.
			On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).
			Return(true, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/1/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"contact_id":  "2",
		})
		rec := httptest.NewRecorder()

		handler.HasSupplierContactRelation(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Verificação concluída", response.Message)

		exists, ok := response.Data.(map[string]any)["exists"].(bool)
		assert.True(t, ok)
		assert.True(t, exists)

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação não existe", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.
			On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).
			Return(false, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/1/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"contact_id":  "2",
		})
		rec := httptest.NewRecorder()

		handler.HasSupplierContactRelation(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		exists, ok := response.Data.(map[string]any)["exists"].(bool)
		assert.True(t, ok)
		assert.False(t, exists)

		mockService.AssertExpectations(t)
	})

	t.Run("error - supplier_id inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/abc/1", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "abc",
			"contact_id":  "1",
		})
		rec := httptest.NewRecorder()

		handler.HasSupplierContactRelation(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - contact_id inválido", func(t *testing.T) {
		_, handler := setupSupplierContact()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/1/xyz", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"contact_id":  "xyz",
		})
		rec := httptest.NewRecorder()

		handler.HasSupplierContactRelation(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - falha interna no service", func(t *testing.T) {
		mockService, handler := setupSupplierContact()

		mockService.
			On("HasSupplierContactRelation", mock.Anything, int64(1), int64(2)).
			Return(false, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/supplier-contact-relations/1/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"contact_id":  "2",
		})
		rec := httptest.NewRecorder()

		handler.HasSupplierContactRelation(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
