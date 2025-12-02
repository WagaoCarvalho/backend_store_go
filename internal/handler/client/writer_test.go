package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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

func TestClientHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

		clientWithVersion := &models.Client{ID: 1, Name: "Cliente Y", Version: 1}
		body, _ := json.Marshal(clientWithVersion)
		req := httptest.NewRequest(http.MethodPut, "/clients/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrZeroVersion).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - corpo inválido", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logAdapter)

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
		handler := NewClientHandler(mockService, logger)

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
		handler := NewClientHandler(mockService, logger)

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
		handler := NewClientHandler(mockService, logger)

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
