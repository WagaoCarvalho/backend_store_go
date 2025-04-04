package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
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

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Data:   users,
		Status: http.StatusOK,
	}

	utils.ToJson(w, response)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["id"]

	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserById(r.Context(), id)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			http.Error(w, `{"status":404, "message":"usuário não encontrado"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"status":500, "message":"Erro interno"}`, http.StatusInternalServerError)
		}
		return
	}

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

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	user, err := h.service.GetUserByEmail(r.Context(), email)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			http.Error(w, `{"status":404, "message":"usuário não encontrado"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"status":500, "message":"Erro interno"}`, http.StatusInternalServerError)
		}
		return
	}

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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		User       models_user.User       `json:"user"`
		CategoryID int64                  `json:"category_id"`
		Address    models_address.Address `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"status":400, "message":"Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), requestData.User, requestData.CategoryID, requestData.Address)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":500, "message":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdUser,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}
	var user models_user.User

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"status":400, "message":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	user.UID = id

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"status":400, "message":"Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), user)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":500, "message":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
		Data:    updatedUser,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUserById(r.Context(), id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "usuário não encontrado") {
			utils.ErrorResponse(w, fmt.Errorf("usuário não encontrado"), http.StatusNotFound)
		} else {
			utils.ErrorResponse(w, fmt.Errorf("erro ao deletar usuário"), http.StatusInternalServerError)
		}
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário deletado com sucesso",
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
