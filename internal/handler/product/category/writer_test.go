package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockService "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryHandler_Create(t *testing.T) {
	baseLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Sucesso - Criar categoria", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		input := dto.ProductCategoryDTO{Name: "Categoria X"}
		expectedModel := dto.ToProductCategoryModel(input)
		expectedModel.ID = int64(1)

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.Name == expectedModel.Name
		})).Return(expectedModel, nil).Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Categoria criada com sucesso", response.Message)

		// Verifica se o ID retornado é int64
		dataBytes, _ := json.Marshal(response.Data)
		var dataDTO dto.ProductCategoryDTO
		json.Unmarshal(dataBytes, &dataDTO)
		assert.Equal(t, expectedModel.ID, *dataDTO.ID)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro JSON inválido", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		// JSON inválido
		body := []byte(`{invalid json`)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro de validação do DTO", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		// Nome vazio deve falhar na validação
		input := dto.ProductCategoryDTO{Name: ""}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro - Categoria já existe", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		input := dto.ProductCategoryDTO{Name: "Categoria X"}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Usar ErrDuplicate que é o erro real do service/repo
		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.Name == input.Name
		})).Return(nil, errMsg.ErrDuplicate).Once()

		h.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "categoria já existe", response.Message)
		assert.Equal(t, http.StatusConflict, response.Status)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico ao criar categoria", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		input := dto.ProductCategoryDTO{Name: "Categoria X"}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.Name == input.Name
		})).Return(nil, errors.New("erro interno")).Once()

		h.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}

func TestProductCategoryHandler_Update(t *testing.T) {
	baseLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	// Na função TestProductCategoryHandler_Update, teste "Sucesso - Update":
	t.Run("Sucesso - Update", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		input := dto.ProductCategoryDTO{Name: "Nova Categoria"}
		expectedModel := dto.ToProductCategoryModel(input)
		expectedModel.ID = id
		updatedModel := &models.ProductCategory{ID: id, Name: "Nova Categoria"} // Linha 33: CORRETO - já é int64

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.ID == id && m.Name == input.Name // Aqui m.ID deve ser int64
		})).Return(nil).Once()

		mockSvc.On("GetByID", mock.Anything, id).Return(updatedModel, nil).Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Categoria atualizada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro - ID inválido", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		req := httptest.NewRequest(http.MethodPut, "/categories/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro - JSON inválido", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer([]byte("{invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro - Categoria não encontrada", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(999)
		input := dto.ProductCategoryDTO{Name: "Nova Categoria"}
		body, _ := json.Marshal(input)

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.ID == id
		})).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPut, "/categories/999", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro - Categoria já existe (conflito)", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		input := dto.ProductCategoryDTO{Name: "Categoria Existente"}
		body, _ := json.Marshal(input)

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.ID == id
		})).Return(errMsg.ErrDuplicate).Once()

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico do service no Update", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		input := dto.ProductCategoryDTO{Name: "Nova Categoria"}
		body, _ := json.Marshal(input)

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.ID == id
		})).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro - Validação do DTO falha (nome vazio)", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		// Nome vazio deve falhar na validação
		input := dto.ProductCategoryDTO{Name: ""}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Service não deve ser chamado
		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("Erro - Validação do DTO falha (nome muito curto)", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		// Nome com 1 caractere deve falhar
		input := dto.ProductCategoryDTO{Name: "A"}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("Erro - Validação do DTO falha (nome muito longo)", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		// Nome com 256 caracteres deve falhar
		longName := strings.Repeat("A", 256)
		input := dto.ProductCategoryDTO{Name: longName}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("Erro - Validação do DTO falha (descrição muito longa)", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		// Descrição com 256 caracteres deve falhar
		longDesc := strings.Repeat("A", 256)
		input := dto.ProductCategoryDTO{
			Name:        "Nome Válido",
			Description: &longDesc,
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("Erro ao buscar categoria atualizada", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		input := dto.ProductCategoryDTO{Name: "Nova Categoria"}
		body, _ := json.Marshal(input)

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(m *models.ProductCategory) bool {
			return m.ID == id && m.Name == input.Name
		})).Return(nil).Once()
		mockSvc.On("GetByID", mock.Anything, id).Return(nil, errors.New("erro ao buscar")).Once()

		req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}

func TestProductCategoryHandler_Delete(t *testing.T) {
	baseLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Sucesso - Delete", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		mockSvc.On("Delete", mock.Anything, id).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro - ID inválido", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		req := httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro - Categoria não encontrada", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(999)
		mockSvc.On("Delete", mock.Anything, id).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		h.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro genérico do service", func(t *testing.T) {
		mockSvc := new(mockService.MockProductCategory)
		h := NewProductCategoryHandler(mockSvc, baseLogger())

		id := int64(1)
		mockSvc.On("Delete", mock.Anything, id).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}
