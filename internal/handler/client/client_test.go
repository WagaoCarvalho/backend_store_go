package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/client"
	dtoClient "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestClientHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Create Client", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
		handler := NewClientHandler(mockService, logAdapter)

		email := "cliente@teste.com"
		cpf := "12345678900"
		inputDTO := &dtoClient.ClientDTO{
			Name:  "Cliente Teste",
			Email: &email,
			CPF:   &cpf,
		}

		expectedModel := dtoClient.ToClientModel(*inputDTO)
		expectedModel.ID = 1

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(m *models.Client) bool {
			return m.Name == inputDTO.Name &&
				m.Email != nil && *m.Email == *inputDTO.Email &&
				m.CPF != nil && *m.CPF == *inputDTO.CPF
		})).Return(expectedModel, nil)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			Status  int                 `json:"status"`
			Message string              `json:"message"`
			Data    dtoClient.ClientDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Cliente criado com sucesso", response.Message)
		assert.NotNil(t, response.Data.ID)
		assert.Equal(t, expectedModel.ID, *response.Data.ID)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid Foreign Key", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
		handler := NewClientHandler(mockService, logAdapter)

		email := "fk@teste.com"
		inputDTO := &dtoClient.ClientDTO{Name: "Cliente FK", Email: &email}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return((*models.Client)(nil), errMsg.ErrInvalidForeignKey)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Duplicate Client", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
		handler := NewClientHandler(mockService, logAdapter)

		email := "duplicado@teste.com"
		inputDTO := &dtoClient.ClientDTO{Name: "Duplicado", Email: &email}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return((*models.Client)(nil), errMsg.ErrDuplicate)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var response utils.DefaultResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		assert.Contains(t, response.Message, "já cadastrado")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
		handler := NewClientHandler(mockService, logAdapter)

		email := "falha@teste.com"
		inputDTO := &dtoClient.ClientDTO{Name: "Falha", Email: &email}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return((*models.Client)(nil), assert.AnError)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
		handler := NewClientHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer([]byte(`{invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestClientHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Get Client by ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
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
			Status  int                 `json:"status"`
			Message string              `json:"message"`
			Data    dtoClient.ClientDTO `json:"data"`
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
		mockService := new(mockClient.MockClientService)
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
		mockService := new(mockClient.MockClientService)
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
		mockService := new(mockClient.MockClientService)
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
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    []dtoClient.ClientDTO `json:"data"`
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
		mockService := new(mockClient.MockClientService)
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
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    []dtoClient.ClientDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Nenhum cliente encontrado", response.Message)
		assert.Empty(t, response.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid param", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientService)
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
		mockService := new(mockClient.MockClientService)
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
