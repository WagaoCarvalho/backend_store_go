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
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
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

// Função auxiliar para extrair ID da URL
func GetIDParam(r *http.Request, key string) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars[key], 10, 64)
}

func ParseErrorResponse(body []byte) (DefaultResponse, error) {
	var resp DefaultResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		// Erro ao fazer o unmarshal, retornando DefaultResponse com status 500
		return DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "Erro ao decodificar resposta",
		}, err
	}
	return resp, nil
}
