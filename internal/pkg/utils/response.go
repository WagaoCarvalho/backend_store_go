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
	Message string      `json:"message,omitempty"`
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

func ToJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "erro ao gerar o JSON", http.StatusInternalServerError)
		log.Println("erro ao converter para JSON:", err)
	}
}

func FromJSON(r io.Reader, target any) error {
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

func GetStringParam(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	val, ok := vars[key]
	if !ok || val == "" {
		return "", fmt.Errorf("missing or empty param: %s", key)
	}
	return val, nil
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

func ParseLimitOffset(r *http.Request) (limit, offset int) {
	query := r.URL.Query()
	limit, _ = strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset, _ = strconv.Atoi(query.Get("offset"))
	return
}

func ParseOrder(r *http.Request) (orderBy, orderDir string) {
	query := r.URL.Query()
	orderBy = query.Get("order_by")
	if orderBy == "" {
		orderBy = "id"
	}
	orderDir = query.Get("order_dir")
	if orderDir == "" {
		orderDir = "asc"
	}
	return
}
