package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DefaultResponse struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

func ErrorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	ToJson(w, struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})

}

func ToJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
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
