package handlers

import (
	"fmt"
	"net/http"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
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
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    users,
	})
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	})
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	user, err := h.service.GetByEmail(r.Context(), email)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	})
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
		Contact    models_contact.Contact `json:"contact"`
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.Create(r.Context(), requestData.User, requestData.CategoryID, requestData.Address, requestData.Contact)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdUser,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		User    models_user.User        `json:"user"`
		Contact *models_contact.Contact `json:"contact,omitempty"`
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	requestData.User.UID = id
	if requestData.Contact != nil {
		requestData.Contact.UserID = &id
	}

	updatedUser, err := h.service.Update(r.Context(), requestData.User, requestData.Contact)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
		Data:    updatedUser,
	})
}

func (h *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário deletado com sucesso",
		Data:    nil,
	})
}
