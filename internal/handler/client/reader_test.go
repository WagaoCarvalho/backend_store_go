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
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Get Client by ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		clientID := int64(10)
		expectedClient := &models.Client{
			ID:      clientID,
			Name:    "Cliente Teste",
			Email:   utils.StrToPtr("teste@cliente.com"),
			CPF:     utils.StrToPtr("123.456.789-00"),
			CNPJ:    nil,
			Version: 1,
			Status:  true,
		}

		mockService.On("GetByID", mock.Anything, clientID).Return(expectedClient, nil)

		req := newRequestWithVars(http.MethodGet, "/clients/"+fmt.Sprint(clientID), nil, map[string]string{"id": fmt.Sprint(clientID)})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int           `json:"status"`
			Message string        `json:"message"`
			Data    dto.ClientDTO `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Cliente encontrado", response.Message)
		assert.Equal(t, expectedClient.ID, *response.Data.ID)
		assert.Equal(t, expectedClient.Name, response.Data.Name)
		assert.Equal(t, *expectedClient.Email, *response.Data.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("Client Not Found", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		clientID := int64(42)
		mockService.On("GetByID", mock.Anything, clientID).Return((*models.Client)(nil), errMsg.ErrNotFound)

		req := newRequestWithVars(http.MethodGet, "/clients/"+fmt.Sprint(clientID), nil, map[string]string{"id": fmt.Sprint(clientID)})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID Param", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		req := newRequestWithVars(http.MethodGet, "/clients/invalid", nil, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestClientHandler_GetByName(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Get Clients by Name", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		name := "Cliente Teste"
		clientModel := &models.Client{
			ID:   1,
			Name: name,
		}

		// agora o mock retorna slice, não apenas 1 objeto
		mockService.On("GetByName", mock.Anything, "Cliente Teste").
			Return([]*models.Client{clientModel}, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/Cliente%20Teste", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Cliente Teste"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int             `json:"status"`
			Message string          `json:"message"`
			Data    []dto.ClientDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Clientes encontrados", response.Message)
		assert.Equal(t, clientModel.ID, *response.Data[0].ID)
		assert.Equal(t, clientModel.Name, response.Data[0].Name)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found - empty list", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		mockService.On("GetByName", mock.Anything, "Inexistente").
			Return([]*models.Client{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/Inexistente", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Inexistente"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode) // agora 200, não 404

		var response struct {
			Status  int             `json:"status"`
			Message string          `json:"message"`
			Data    []dto.ClientDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Nenhum cliente encontrado", response.Message)
		assert.Empty(t, response.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid param", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/", nil)
		req = mux.SetURLVars(req, map[string]string{"name": ""})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		mockService.On("GetByName", mock.Anything, "Erro").
			Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/Erro", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Erro"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

}

func TestClientHandler_GetVersionByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/clients/invalid/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		mockService.
			On("GetVersionByID", mock.Anything, int64(123)).
			Return(0, errors.New("erro no banco"))

		req := httptest.NewRequest(http.MethodGet, "/clients/123/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - versão encontrada", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		expectedVersion := 5

		mockService.
			On("GetVersionByID", mock.Anything, int64(123)).
			Return(expectedVersion, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/123/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Versão do cliente recuperada com sucesso", resp.Message)

		dataMap, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(123), dataMap["client_id"]) // JSON converte int64 para float64
		assert.Equal(t, float64(expectedVersion), dataMap["version"])

		mockService.AssertExpectations(t)
	})
}

func TestClientHandler_ClientExists(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)

	setup := func() (*mockClient.MockClient, *clientHandler) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, loggerAdapter)
		return mockService, handler
	}

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/clients/invalid/exists", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.ClientExists(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna erro", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)
		mockService.On("ClientExists", mock.Anything, clientID).Return(false, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/clients/1/exists", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ClientExists(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - cliente existe", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)
		mockService.On("ClientExists", mock.Anything, clientID).Return(true, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/clients/1/exists", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ClientExists(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Status)
		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(1), data["client_id"]) // JSON unmarshals numbers como float64
		assert.Equal(t, true, data["exists"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - cliente não existe", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(2)
		mockService.On("ClientExists", mock.Anything, clientID).Return(false, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/clients/2/exists", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.ClientExists(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Status)
		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(2), data["client_id"])
		assert.Equal(t, false, data["exists"])

		mockService.AssertExpectations(t)
	})
}
