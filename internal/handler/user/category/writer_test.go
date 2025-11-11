package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockUserCat "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/category"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCategoryHandler_Create(t *testing.T) {
	mockSvc := new(mockUserCat.MockUserCategory)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUserCategoryHandler(mockSvc, logger)

	t.Run("Success", func(t *testing.T) {
		categoryDTO := dto.UserCategoryDTO{Name: "Nova"}
		categoryModel := dto.ToUserCategoryModel(categoryDTO)

		mockSvc.On("Create", mock.Anything, categoryModel).Return(categoryModel, nil)

		body, _ := json.Marshal(categoryDTO)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		var result dto.UserCategoryDTO
		itemBytes, _ := json.Marshal(response.Data)
		err = json.Unmarshal(itemBytes, &result)
		require.NoError(t, err)

		assert.Equal(t, categoryDTO.Name, result.Name)
		assert.Equal(t, "Categoria criada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer([]byte("{invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "erro ao decodificar JSON")
	})

	t.Run("ServiceError", func(t *testing.T) {
		inputDTO := dto.UserCategoryDTO{Name: "Erro"}
		inputModel := dto.ToUserCategoryModel(inputDTO)

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.Name == inputModel.Name
		})).Return(nil, errors.New("erro ao criar categoria"))

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, "erro ao criar categoria", resp.Message)

		mockSvc.AssertExpectations(t)
	})
}

func TestUserCategoryHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		categoryDTO := dto.UserCategoryDTO{Name: "Atualizada"}
		categoryModel := dto.ToUserCategoryModel(categoryDTO)
		categoryModel.ID = 1

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.ID == 1 && c.Name == "Atualizada"
		})).Return(nil)

		body, _ := json.Marshal(categoryDTO)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/1", bytes.NewBuffer(body)), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Categoria atualizada com sucesso", response.Message)
		assert.Equal(t, http.StatusOK, response.Status)

		mockSvc.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		categoryDTO := dto.UserCategoryDTO{Name: "Inexistente"}

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.ID == 999
		})).Return(errMsg.ErrNotFound)

		body, _ := json.Marshal(categoryDTO)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/999", bytes.NewBuffer(body)), map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "categoria não encontrada", resp.Message)
		assert.Equal(t, http.StatusNotFound, resp.Status)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(mockUserCat.MockUserCategory), logger)

		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/abc", bytes.NewBuffer([]byte("{}"))), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "ID inválido")
		assert.Equal(t, http.StatusBadRequest, resp.Status)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(mockUserCat.MockUserCategory), logger)

		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/1", bytes.NewBuffer([]byte("{invalid"))), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "erro ao decodificar JSON")
		assert.Equal(t, http.StatusBadRequest, resp.Status)
	})

	t.Run("UpdateError", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		categoryDTO := dto.UserCategoryDTO{Name: "Falha"}
		categoryModel := dto.ToUserCategoryModel(categoryDTO)
		categoryModel.ID = 2

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.ID == 2
		})).Return(errors.New("erro ao atualizar"))

		body, _ := json.Marshal(categoryDTO)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/2", bytes.NewBuffer(body)), map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "erro ao atualizar categoria")
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ZeroID", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		categoryDTO := dto.UserCategoryDTO{Name: "SemID"}
		categoryModel := dto.ToUserCategoryModel(categoryDTO)
		categoryModel.ID = 10

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.ID == 10
		})).Return(errMsg.ErrZeroID)

		body, _ := json.Marshal(categoryDTO)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/10", bytes.NewBuffer(body)), map[string]string{"id": "10"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidData", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		categoryDTO := dto.UserCategoryDTO{Name: ""}
		categoryModel := dto.ToUserCategoryModel(categoryDTO)
		categoryModel.ID = 5

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.ID == 5
		})).Return(errMsg.ErrInvalidData)

		body, _ := json.Marshal(categoryDTO)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/5", bytes.NewBuffer(body)), map[string]string{"id": "5"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "dados inválidos")

		mockSvc.AssertExpectations(t)
	})

}

func TestUserCategoryHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(mockUserCat.MockUserCategory), logger)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/abc", nil), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserCat.MockUserCategory)
		handler := NewUserCategoryHandler(mockSvc, logger)

		mockSvc.On("Delete", mock.Anything, int64(10)).Return(errors.New("erro ao deletar"))

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/10", nil), map[string]string{"id": "10"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, "erro ao deletar", resp.Message)

		mockSvc.AssertExpectations(t)
	})
}
