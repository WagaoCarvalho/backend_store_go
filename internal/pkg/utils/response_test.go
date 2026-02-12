package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	ErrorResponse(rr, errors.New("erro teste"), http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp DefaultResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "erro teste", resp.Message)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
}

func TestErrorResponse_NilError(t *testing.T) {
	rr := httptest.NewRecorder()
	ErrorResponse(rr, nil, http.StatusInternalServerError)

	var resp DefaultResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "", resp.Message)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
}

func TestToJSON_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	ToJSON(rr, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result map[string]string
	err := json.NewDecoder(rr.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestFromJSON_Success(t *testing.T) {
	input := `{"name":"Test"}`
	var result map[string]string

	err := FromJSON(bytes.NewBufferString(input), &result)
	assert.NoError(t, err)
	assert.Equal(t, "Test", result["name"])
}

func TestFromJSON_Invalid(t *testing.T) {
	input := `invalid json`
	var result map[string]string

	err := FromJSON(bytes.NewBufferString(input), &result)
	assert.Error(t, err)
}

func TestGetIDParam_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})

	id, err := GetIDParam(req, "id")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
}

func TestGetIDParam_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	_, err := GetIDParam(req, "id")
	assert.Error(t, err)
}

func TestGetIDParam_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})

	_, err := GetIDParam(req, "id")
	assert.Error(t, err)
}

func TestGetStringParam_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/john", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "john"})

	val, err := GetStringParam(req, "name")
	assert.NoError(t, err)
	assert.Equal(t, "john", val)
}

func TestGetStringParam_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	_, err := GetStringParam(req, "name")
	assert.Error(t, err)
}

func TestGetPaginationParams_Default(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?x=1", nil)

	limit, offset := GetPaginationParams(req)

	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, offset)
}

func TestGetPaginationParams_Custom(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=20&offset=5", nil)

	limit, offset := GetPaginationParams(req)

	assert.Equal(t, 20, limit)
	assert.Equal(t, 5, offset)
}

func TestGetPaginationParams_LimitZeroOrNegative(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=0", nil)

	limit, _ := GetPaginationParams(req)

	assert.Equal(t, 10, limit)
}

func TestGetPaginationParams_LimitAboveMax(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=200", nil)

	limit, _ := GetPaginationParams(req)

	assert.Equal(t, 100, limit)
}

func TestGetPaginationParams_InvalidOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?offset=-10", nil)

	_, offset := GetPaginationParams(req)

	assert.Equal(t, 0, offset)
}

func TestParseErrorResponse_Success(t *testing.T) {
	body := []byte(`{"status":404,"message":"Não encontrado"}`)

	resp, err := ParseErrorResponse(body)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Status)
	assert.Equal(t, "Não encontrado", resp.Message)
}

func TestParseErrorResponse_InvalidJSON(t *testing.T) {
	body := []byte(`{{{{`)

	resp, err := ParseErrorResponse(body)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
}

func TestParseOrder_Default(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	orderBy, orderDir := ParseOrder(req)

	assert.Equal(t, "id", orderBy)
	assert.Equal(t, "asc", orderDir)
}

func TestParseOrder_Custom(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?order_by=name&order_dir=desc", nil)

	orderBy, orderDir := ParseOrder(req)

	assert.Equal(t, "name", orderBy)
	assert.Equal(t, "desc", orderDir)
}

// Writer que simula erro
type failingWriter struct {
	header http.Header
}

func (f *failingWriter) Header() http.Header {
	if f.header == nil {
		f.header = make(http.Header)
	}
	return f.header
}

func (f *failingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("erro simulado")
}

func (f *failingWriter) WriteHeader(_ int) {}

func TestErrorResponse_EncodeError(t *testing.T) {
	w := &failingWriter{}
	ErrorResponse(w, errors.New("erro"), http.StatusBadRequest)
	assert.True(t, true)
}

func TestToJSON_EncodeError(t *testing.T) {
	log.SetFlags(0)

	w := &failingWriter{}
	dado := map[string]string{"k": "v"}

	ToJSON(w, http.StatusOK, dado)

	assert.True(t, true)
}

func TestParseLimitOffset_Default(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	limit, offset := ParseLimitOffset(req)

	cfg := config.LoadPaginationConfig()

	assert.Equal(t, cfg.DefaultLimit, limit)
	assert.Equal(t, cfg.DefaultOffset, offset)
}

func TestParseLimitOffset_CustomValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=25&offset=5", nil)

	limit, offset := ParseLimitOffset(req)

	assert.Equal(t, 25, limit)
	assert.Equal(t, 5, offset)
}

func TestParseLimitOffset_InvalidLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=0", nil)

	limit, _ := ParseLimitOffset(req)

	cfg := config.LoadPaginationConfig()

	// Deve manter default porque limit <= 0
	assert.Equal(t, cfg.DefaultLimit, limit)
}

func TestParseLimitOffset_InvalidOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?offset=-10", nil)

	_, offset := ParseLimitOffset(req)

	cfg := config.LoadPaginationConfig()

	// Deve manter default porque offset < 0
	assert.Equal(t, cfg.DefaultOffset, offset)
}

func TestParseLimitOffset_LimitAboveMax(t *testing.T) {
	cfg := config.LoadPaginationConfig()

	req := httptest.NewRequest(http.MethodGet, "/?limit=9999", nil)

	limit, _ := ParseLimitOffset(req)

	assert.Equal(t, cfg.MaxLimit, limit)
}

func TestParseLimitOffset_InvalidNumbers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=abc&offset=xyz", nil)

	limit, offset := ParseLimitOffset(req)

	cfg := config.LoadPaginationConfig()

	// Deve manter defaults
	assert.Equal(t, cfg.DefaultLimit, limit)
	assert.Equal(t, cfg.DefaultOffset, offset)
}
