package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockService "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryHandler_GetByID(t *testing.T) {
	baseLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Sucesso - GetByID", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategory(mockSvc, baseLogger())

		id := int64(1)
		expectedModel := &models.ProductCategory{ID: uint(id), Name: "Categoria X"}

		mockSvc.On("GetByID", mock.Anything, id).Return(expectedModel, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int                    `json:"status"`
			Message string                 `json:"message"`
			Data    dto.ProductCategoryDTO `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedModel.ID, *response.Data.ID)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro - ID inválido", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategory(mockSvc, baseLogger())

		req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro - Categoria não encontrada", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategory(mockSvc, baseLogger())

		id := int64(999)
		mockSvc.On("GetByID", mock.Anything, id).Return(nil, errors.New("categoria não encontrada")).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico do service", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategory(mockSvc, baseLogger())

		id := int64(1)
		mockSvc.On("GetByID", mock.Anything, id).Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}

func TestProductCategoryHandler_GetAll(t *testing.T) {
	baseLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}
	t.Run("Sucesso - GetAll", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategory(mockSvc, baseLogger())

		expectedModels := []*models.ProductCategory{
			{ID: 1, Name: "Categoria 1"},
			{ID: 2, Name: "Categoria 2"},
		}

		mockSvc.On("GetAll", mock.Anything).Return(expectedModels, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int                      `json:"status"`
			Message string                   `json:"message"`
			Data    []dto.ProductCategoryDTO `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, len(expectedModels))

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico do service - GetAll", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategory(mockSvc, baseLogger())

		mockSvc.On("GetAll", mock.Anything).Return([]*models.ProductCategory(nil), errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

}
