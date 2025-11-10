package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
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
		handler := NewClient(mockService, logger)

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

	t.Run("sucesso - clientes listados", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		clients := []*models.Client{
			{ID: 1, Name: "Cliente 1"},
			{ID: 2, Name: "Cliente 2"},
		}

		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return(clients, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?limit=5&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Clientes listados com sucesso", resp.Message)

		var data []dto.ClientDTO
		dataBytes, err := json.Marshal(resp.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(dataBytes, &data)
		assert.NoError(t, err)

		assert.Len(t, data, 2)
		assert.Equal(t, int64(1), *data[0].ID)
		assert.Equal(t, "Cliente 1", data[0].Name)
		assert.Equal(t, int64(2), *data[1].ID)
		assert.Equal(t, "Cliente 2", data[1].Name)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtrando por name e status", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		clients := []*models.Client{
			{ID: 3, Name: "Empresa XPTO"},
		}

		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return(clients, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/filter?name=Empresa&status=true&limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data, ok := resp.Data.([]interface{})
		assert.True(t, ok, "Data should be a slice of interfaces")
		assert.Equal(t, 1, len(data))

		mockService.AssertExpectations(t)
	})

	t.Run("erro - filtro inválido", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		// Configure o mock para retornar erro de filtro inválido
		mockService.
			On("GetAll", mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("%w: intervalo de criação inválido", errMsg.ErrInvalidFilter))

		// Cria request com filtro inválido (datas invertidas)
		req := httptest.NewRequest(http.MethodGet, "/clients/filter?created_from=2025-12-01&created_to=2025-01-01", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

}
