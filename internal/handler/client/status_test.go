package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientHandler_Disable(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockClient.MockClient, *clientHandler) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
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

	t.Run("Invalid Method", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/clients/disable/1", nil)
		w := httptest.NewRecorder()

		handler.Disable(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
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

	t.Run("ErrVersionConflict", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)
		mockService.On("Disable", mock.Anything, clientID).Return(errMsg.ErrZeroVersion).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Other Error", func(t *testing.T) {
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
}

func TestClientHandler_Enable(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockClient.MockClient, *clientHandler) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
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

	t.Run("Invalid Method", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPost, "/clients/enable/1", nil)
		w := httptest.NewRecorder()

		handler.Enable(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
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

	t.Run("ErrVersionConflict", func(t *testing.T) {
		mockService, handler := setup()
		clientID := int64(1)
		mockService.On("Enable", mock.Anything, clientID).Return(errMsg.ErrZeroVersion).Once()

		req := httptest.NewRequest(http.MethodPatch, "/clients/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Other Error", func(t *testing.T) {
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
