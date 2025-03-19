package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WagaoCarvalho/backend_store_go/internal/services"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Chama o serviço para obter os usuários
	users, err := h.service.GetUsers(ctx)
	if err != nil {
		// Erro ao buscar usuários, retorna uma resposta de erro
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	// Cria a resposta padrão
	response := utils.DefaultResponse{
		Data:   users,
		Status: http.StatusOK,
	}

	// Converte a resposta para JSON e envia para o cliente
	utils.ToJson(w, response)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["id"]

	// Convertendo o ID para int64
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Buscando o usuário
	user, err := h.service.GetUserById(r.Context(), id)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			http.Error(w, `{"status":404, "message":"usuário não encontrado"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"status":500, "message":"Erro interno"}`, http.StatusInternalServerError)
		}
		return
	}

	// Enviar resposta com dados do usuário
	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao gerar resposta", http.StatusInternalServerError)
	}
}
