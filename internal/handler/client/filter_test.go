package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientHandler_GetAll(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - falha no serviço", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - filtro inválido", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter)

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?limit=-1", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista de clientes", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		mockClients := []*dto.ClientDTO{
			{
				ID:          utils.Int64Ptr(1),
				Name:        "João da Silva",
				Email:       utils.StrToPtr("joao@teste.com"),
				CPF:         utils.StrToPtr("12345678900"),
				Description: "Cliente teste",
				Status:      true,
				Version:     1,
			},
			{
				ID:          utils.Int64Ptr(2),
				Name:        "Maria Souza",
				Email:       utils.StrToPtr("maria@teste.com"),
				CPF:         utils.StrToPtr("98765432100"),
				Description: "Cliente teste 2",
				Status:      false,
				Version:     1,
			},
		}

		// Converta para model.Client, já que o serviço retorna o model, não o DTO
		modelClients := []*model.Client{}
		for _, c := range mockClients {
			modelClients = append(modelClients, &model.Client{
				ID:          int64(*c.ID),
				Name:        c.Name,
				Email:       c.Email,
				CPF:         c.CPF,
				Description: c.Description,
				Status:      c.Status,
				Version:     c.Version,
			})
		}

		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return(modelClients, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Clientes listados com sucesso", resp.Message)

		data, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(2), data["total"]) // JSON unmarshal converte para float64
		assert.NotEmpty(t, data["items"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista vazia", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return([]*model.Client{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?limit=5&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)

		data, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(0), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com status true", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		// Mock de retorno
		mockClients := []*model.Client{
			{
				ID:     1,
				Name:   "João da Silva",
				Email:  utils.StrToPtr("joao@teste.com"),
				Status: true,
			},
		}

		mockService.
			On("GetAll", mock.Anything, mock.MatchedBy(func(f *model.ClientFilter) bool {
				return f.Status != nil && *f.Status == true
			})).
			Return(mockClients, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?status=true", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(1), data["total"]) // JSON converte int → float64
		items := data["items"].([]interface{})
		assert.Equal(t, "João da Silva", items[0].(map[string]interface{})["name"])

		mockService.AssertExpectations(t)
	})

}
