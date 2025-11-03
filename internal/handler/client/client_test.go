package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
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
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		email := "cliente@teste.com"
		cpf := "12345678900"
		inputDTO := &dto.ClientDTO{
			Name:  "Cliente Teste",
			Email: &email,
			CPF:   &cpf,
		}

		expectedModel := dto.ToClientModel(*inputDTO)
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
			Status  int           `json:"status"`
			Message string        `json:"message"`
			Data    dto.ClientDTO `json:"data"`
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
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		email := "fk@teste.com"
		inputDTO := &dto.ClientDTO{Name: "Cliente FK", Email: &email}

		mockService.On("Create", mock.Anything, mock.Anything).
			Return((*models.Client)(nil), errMsg.ErrDBInvalidForeignKey)

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
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		email := "duplicado@teste.com"
		inputDTO := &dto.ClientDTO{Name: "Duplicado", Email: &email}

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
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		email := "falha@teste.com"
		inputDTO := &dto.ClientDTO{Name: "Falha", Email: &email}

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
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

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
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logAdapter)

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
		handler := NewClient(mockService, logger)

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
		handler := NewClient(mockService, logger)

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
		handler := NewClient(mockService, logger)

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

func TestClientHandler_GetAll(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - falha no serviço", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		mockService.
			On("GetAll", mock.Anything, 10, 0).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/clients?limit=10&offset=0", nil)
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
			On("GetAll", mock.Anything, 5, 0).
			Return(clients, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients?limit=5&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Clientes listados com sucesso", resp.Message)

		// Reinterpreta Data como slice de DTOs
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
}

func TestClientHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/clients/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - dados inválidos (ErrInvalidData)", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		// Simula um client com dados inválidos
		invalidClient := &models.Client{ID: 1, Name: ""} // Name obrigatório, por exemplo

		body, _ := json.Marshal(invalidClient)
		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrInvalidData).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - ID zero (ErrZeroID)", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		clientWithZeroID := &models.Client{ID: 0, Name: "Cliente X"}
		body, _ := json.Marshal(clientWithZeroID)
		req := httptest.NewRequest(http.MethodPut, "/clients/0", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrZeroID).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - conflito de versão (ErrVersionConflict)", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		clientWithVersion := &models.Client{ID: 1, Name: "Cliente Y", Version: 1}
		body, _ := json.Marshal(clientWithVersion)
		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrVersionConflict).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - corpo inválido", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewBuffer([]byte("{invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - service retorna ErrNotFound", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		inputDTO := &dto.ClientDTO{Name: "Cliente Teste"}
		body, _ := json.Marshal(inputDTO)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Client) bool {
			return c.ID == 1 && c.Name == "Cliente Teste"
		})).Return(errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - service retorna ErrDuplicate", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		inputDTO := &dto.ClientDTO{Name: "Cliente Teste"}
		body, _ := json.Marshal(inputDTO)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Client) bool {
			return c.ID == 1 && c.Name == "Cliente Teste"
		})).Return(errMsg.ErrDuplicate)

		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - update", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		email := "cliente@teste.com"
		cpf := "12345678900"
		inputDTO := &dto.ClientDTO{
			Name:  "Cliente Teste Atualizado",
			Email: &email,
			CPF:   &cpf,
		}

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Client) bool {
			return m.Name == inputDTO.Name &&
				m.Email != nil && *m.Email == *inputDTO.Email &&
				m.CPF != nil && *m.CPF == *inputDTO.CPF &&
				m.ID == 1
		})).Return(nil)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

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
		assert.Equal(t, "Cliente atualizado com sucesso", response.Message)
		assert.Equal(t, inputDTO.Name, response.Data.Name)
		assert.Equal(t, *inputDTO.Email, *response.Data.Email)
		assert.Equal(t, *inputDTO.CPF, *response.Data.CPF)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - service retorna outro erro", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)

		uid := "1"
		reqBody := &dto.ClientDTO{Name: "Cliente Teste"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": uid})
		rec := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Client) bool {
			return c.ID == 1
		})).Return(errors.New("erro inesperado"))

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var resp utils.DefaultResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

}

func TestClientHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/clients/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna erro", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		mockService.On("Delete", mock.Anything, int64(1)).Return(errors.New("erro no banco"))

		req := httptest.NewRequest(http.MethodDelete, "/clients/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - cliente deletado", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logger)

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/clients/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())

		mockService.AssertExpectations(t)
	})
}

func TestClientHandler_DisableEnable(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockClient.MockClient, *Client) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Disable - Success", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)

		mockService.On("Disable", mock.Anything, clientID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Enable - Success", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)

		mockService.On("Enable", mock.Anything, clientID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Disable - Invalid Method", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/clients/disable/1", nil)
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Enable - Invalid Method", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPost, "/clients/enable/1", nil)
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Disable - ID inválido", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Enable - ID inválido", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disable - ErrNotFound", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)

		mockService.On("Disable", mock.Anything, clientID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Enable - ErrNotFound", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)

		mockService.On("Enable", mock.Anything, clientID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Disable - ErrVersionConflict", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)

		mockService.On("Disable", mock.Anything, clientID).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Enable - ErrVersionConflict", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)

		mockService.On("Enable", mock.Anything, clientID).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Disable - Other Error", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)
		mockService.On("Disable", mock.Anything, clientID).Return(errors.New("other error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Enable - Other Error", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)
		mockService.On("Enable", mock.Anything, clientID).Return(errors.New("other error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestClientHandler_ClientExists(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)

	setup := func() (*mockClient.MockClient, *Client) {
		mockService := new(mockClient.MockClient)
		handler := NewClient(mockService, loggerAdapter)
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
