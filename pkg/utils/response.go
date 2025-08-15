package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DefaultResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"` // <- Usado tanto para sucesso quanto erro
	Status  int         `json:"status"`
}

func ErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var message string
	if err != nil {
		message = err.Error()
	}

	response := DefaultResponse{
		Status:  statusCode,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao codificar a resposta", http.StatusInternalServerError)
	}
}

func ToJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "erro ao gerar o JSON", http.StatusInternalServerError)
		log.Println("erro ao converter para JSON:", err)
	}
}

func FromJson(r io.Reader, target interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}
	return nil
}

func GetIDParam(r *http.Request, key string) (int64, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[key]
	if !ok || idStr == "" {
		return 0, fmt.Errorf("missing or empty id param: %s", key)
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %s", idStr)
	}
	return id, nil
}

func ParseErrorResponse(body []byte) (DefaultResponse, error) {
	var resp DefaultResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "Erro ao decodificar resposta",
		}, err
	}
	return resp, nil
}
