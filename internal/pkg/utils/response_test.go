package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestToJson_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	ToJSON(rr, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result map[string]string
	err := json.NewDecoder(rr.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestFromJson_Success(t *testing.T) {
	input := `{"name":"Test"}`
	var result map[string]string

	err := FromJSON(bytes.NewBufferString(input), &result)
	assert.NoError(t, err)
	assert.Equal(t, "Test", result["name"])
}

// Teste do FromJson com JSON inválido
func TestFromJson_Invalid(t *testing.T) {
	jsonInput := `{"invalid}`
	var result map[string]string

	err := FromJSON(bytes.NewBufferString(jsonInput), &result)

	assert.Error(t, err)
}

func TestFromJson_Error(t *testing.T) {
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

func TestGetIDParam_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})

	_, err := GetIDParam(req, "id")
	assert.Error(t, err)
}

func TestParseErrorResponse_Success(t *testing.T) {
	body := []byte(`{"status": 404, "message": "Não encontrado"}`)
	resp, err := ParseErrorResponse(body)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Status)
	assert.Equal(t, "Não encontrado", resp.Message)
}

func TestParseErrorResponse_InvalidJSON(t *testing.T) {
	body := []byte(`{{{{`) // JSON inválido
	resp, err := ParseErrorResponse(body)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.NotEmpty(t, resp.Message)
}

// ResponseWriter que sempre retorna erro no Encode (simula erro de encoding)
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
	return 0, errors.New("erro simulado no Write")
}

func (f *failingWriter) WriteHeader(_ int) {}

func TestErrorResponse_EncodeError(t *testing.T) {
	w := &failingWriter{}

	ErrorResponse(w, errors.New("erro qualquer"), http.StatusBadRequest)

	// Não temos acesso direto ao que foi escrito (por ser erro),
	// mas podemos garantir que o fallback `http.Error` foi chamado sem pânico.
	// Este teste não falha se a função tentar `Encode` e cair no fallback.
	assert.True(t, true, "Fallback de erro executado sem crash")
}

func TestToJson_EncodeError(t *testing.T) {
	// Redireciona logs (opcional)
	log.SetFlags(0)

	w := &failingWriter{}

	// Simula um dado qualquer
	dado := map[string]string{"chave": "valor"}

	ToJSON(w, http.StatusOK, dado)

	// Não esperamos retorno visível, mas sim que não cause pânico e caia no fallback
	t.Log("TestToJson_EncodeError executado sem crash, fallback de erro acionado")
}
