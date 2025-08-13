package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type response struct {
	Message string `json:"message"`
}

func TestGetHome(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h := http.HandlerFunc(GetHome)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperado status %d, mas recebeu %d", http.StatusOK, rr.Code)
	}

	var res response
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("Erro ao fazer parse do JSON: %v", err)
	}

	expectedMessage := "Go RESTful Api backend_store"
	if res.Message != expectedMessage {
		t.Errorf("Esperado '%s', mas recebeu '%s'", expectedMessage, res.Message)
	}
}
