package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressHandler_Enable(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Enable Address", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Enable", mock.Anything, id).Return(nil)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/invalid/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Error - Enable Service returns ErrNotFound", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Enable", mock.Anything, id).Return(errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Error - Enable Service returns generic error", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Enable", mock.Anything, id).Return(errors.New("db error"))

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Error - Enable Method Not Allowed", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/1/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Error - Enable Invalid ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/abc/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Error - Enable Service returns ErrID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Enable", mock.Anything, id).Return(errMsg.ErrZeroID)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_Disable(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Disable Address", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Disable", mock.Anything, id).Return(nil)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/invalid/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Error - Disable Service returns ErrNotFound", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Disable", mock.Anything, id).Return(errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Error - Disable Service returns generic error", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Disable", mock.Anything, id).Return(errors.New("db error"))

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Error - Disable Method Not Allowed", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/1/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Error - Disable Invalid ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/abc/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Error - Disable Service returns ErrID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddress)
		handler := NewAddress(mockService, logAdapter)

		id := int64(1)
		mockService.On("Disable", mock.Anything, id).Return(errMsg.ErrZeroID)

		req := httptest.NewRequest(http.MethodPatch, "/addresses/1/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
