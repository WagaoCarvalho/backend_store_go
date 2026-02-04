package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/contact"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestLogger() *logger.LogAdapter {
	base := logrus.New()
	base.Out = &bytes.Buffer{}
	return logger.NewLoggerAdapter(base)
}

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

/* ========================= CREATE ========================= */

func TestContactHandler_Create(t *testing.T) {
	logAdapter := newTestLogger()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		dto := dtoContact.ContactDTO{
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123",
		}

		model := dtoContact.ToContactModel(dto)
		mockSvc.On("Create", mock.Anything, model).
			Return((*models.Contact)(model), nil)

		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader("{invalid"))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ErrInvalidData", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		dto := dtoContact.ContactDTO{ContactName: ""}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Create", mock.Anything, model).
			Return(nil, errMsg.ErrInvalidData)

		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ErrDuplicate", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		dto := dtoContact.ContactDTO{ContactName: "Contato"}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Create", mock.Anything, model).
			Return(nil, errMsg.ErrDuplicate)

		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		dto := dtoContact.ContactDTO{ContactName: "Contato"}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Create", mock.Anything, model).
			Return(nil, errMsg.ErrNotFound)

		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		dto := dtoContact.ContactDTO{ContactName: "Contato"}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Create", mock.Anything, model).
			Return(nil, errors.New("erro inesperado"))

		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

/* ========================= UPDATE ========================= */

/* ========================= UPDATE ========================= */

/* ========================= UPDATE ========================= */

func TestContactHandler_Update(t *testing.T) {
	logAdapter := newTestLogger()

	toReader := func(v any) io.Reader {
		b, _ := json.Marshal(v)
		return bytes.NewReader(b)
	}

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		id := int64(1)
		dto := dtoContact.ContactDTO{
			ID:          &id,
			ContactName: "Nome",
		}

		model := dtoContact.ToContactModel(dto)
		mockSvc.On("Update", mock.Anything, model).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id parameter", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/contacts/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("zero id", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		id := int64(0)
		dto := dtoContact.ContactDTO{ID: &id}

		req := httptest.NewRequest(http.MethodPut, "/contacts/0", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("invalid json body", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", strings.NewReader("{invalid json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Update")

		// Verificar se retornou um JSON válido com status de erro
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "A resposta deve ser um JSON válido")
		assert.Equal(t, float64(400), response["status"])
		assert.Contains(t, response, "message")
	})

	t.Run("empty body", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Update")
	})

	t.Run("ErrInvalidData", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Update", mock.Anything, model).
			Return(errMsg.ErrInvalidData)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Update", mock.Anything, model).
			Return(errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ErrDuplicate", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		id := int64(1)
		dto := dtoContact.ContactDTO{
			ID:          &id,
			ContactName: "Nome Duplicado",
		}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Update", mock.Anything, model).
			Return(errMsg.ErrDuplicate)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id}
		model := dtoContact.ToContactModel(dto)

		mockSvc.On("Update", mock.Anything, model).
			Return(errors.New("erro inesperado"))

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

/* ========================= DELETE ========================= */

func TestContactHandler_Delete(t *testing.T) {
	logAdapter := newTestLogger()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := newRequestWithVars("DELETE", "/contacts/1", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("DELETE", "/contacts/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Delete")
	})

	t.Run("zero id", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("DELETE", "/contacts/0", nil, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Delete")
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(99)).
			Return(errMsg.ErrNotFound)

		req := newRequestWithVars("DELETE", "/contacts/99", nil, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockSvc := new(mockContact.MockContact)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(2)).
			Return(errors.New("erro inesperado"))

		req := newRequestWithVars("DELETE", "/contacts/2", nil, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
