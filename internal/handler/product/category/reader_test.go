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
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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

	t.Run("Sucesso", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		expectedModel := &models.ProductCategory{ID: id, Name: "Categoria X"}

		mockSvc.On("GetByID", mock.Anything, id).Return(expectedModel, nil)

		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ID inválido", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Categoria não encontrada", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(999)
		mockSvc.On("GetByID", mock.Anything, id).Return(nil, errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ID zero", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(0)
		mockSvc.On("GetByID", mock.Anything, id).Return(nil, errMsg.ErrZeroID)

		req := httptest.NewRequest(http.MethodGet, "/categories/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		mockSvc.On("GetByID", mock.Anything, id).Return(nil, errors.New("erro"))

		req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestProductCategoryHandler_GetAll(t *testing.T) {
	baseLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Sucesso - GetAll com dados", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		expectedModels := []*models.ProductCategory{
			{ID: int64(1), Name: "Categoria 1", Description: "Desc 1"},
			{ID: int64(2), Name: "Categoria 2", Description: "Desc 2"},
			{ID: int64(3), Name: "Categoria 3", Description: ""},
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
		assert.Equal(t, "Categorias recuperadas com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Sucesso - GetAll sem dados (lista vazia)", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		expectedModels := []*models.ProductCategory{}

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
		assert.Empty(t, response.Data)
		assert.Equal(t, "Categorias recuperadas com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico do service - GetAll", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		dbError := errors.New("erro de conexão com banco")
		mockSvc.On("GetAll", mock.Anything).Return([]*models.ProductCategory(nil), dbError).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "erro ao buscar categorias")

		mockSvc.AssertExpectations(t)
	})

}
