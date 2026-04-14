package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_Enable(t *testing.T) {
	t.Run("successfully enable supplier", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		supplierID := int64(1)
		requestBody := map[string]interface{}{
			"version": 2,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("GetByID", mock.Anything, supplierID).
			Return(&models.Supplier{
				ID:      supplierID,
				Status:  false,
				Version: 2,
			}, nil).Once()

		mockService.On("Enable", mock.Anything, supplierID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return method not allowed for wrong HTTP method", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodGet, "/suppliers/1/enable", nil)
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("return bad request when id is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/abc/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return bad request when version is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":0}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return bad request when payload is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte("invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return not found when get by id returns ErrNotFound", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when get by id fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return conflict when version mismatch", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  false,
				Version: 1,
			}, nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":2}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return not found when enable service returns ErrNotFound", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  false,
				Version: 1,
			}, nil).Once()

		mockService.On("Enable", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when enable service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  false,
				Version: 1,
			}, nil).Once()

		mockService.On("Enable", mock.Anything, int64(1)).Return(errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_Disable(t *testing.T) {
	t.Run("successfully disable supplier", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		supplierID := int64(1)
		requestBody := map[string]interface{}{
			"version": 2,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("GetByID", mock.Anything, supplierID).
			Return(&models.Supplier{
				ID:      supplierID,
				Status:  true,
				Version: 2,
			}, nil).Once()

		mockService.On("Disable", mock.Anything, supplierID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return method not allowed for wrong HTTP method", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPost, "/suppliers/1/disable", nil)
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("return bad request when id is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/abc/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return bad request when version is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":-1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return bad request when payload is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte("invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return not found when get by id returns ErrNotFound", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when get by id fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return conflict when version mismatch", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  true,
				Version: 1,
			}, nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":2}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return not found when disable service returns ErrNotFound", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  true,
				Version: 1,
			}, nil).Once()

		mockService.On("Disable", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when disable service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  true,
				Version: 1,
			}, nil).Once()

		mockService.On("Disable", mock.Anything, int64(1)).Return(errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
