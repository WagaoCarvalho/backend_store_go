package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	user_category_relations_models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
)

type UserCategoryRelationHandler struct {
	service services.UserCategoryRelationServices
}

func NewUserCategoryRelationHandler(service services.UserCategoryRelationServices) *UserCategoryRelationHandler {
	return &UserCategoryRelationHandler{service: service}
}

func (h *UserCategoryRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var relation *user_category_relations_models.UserCategoryRelations
	if err := utils.FromJson(r.Body, &relation); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	created, wasCreated, err := h.service.Create(ctx, relation.UserID, relation.CategoryID)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if errors.Is(err, repositories.ErrInvalidForeignKey) {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var status int
	var message string

	if wasCreated {
		status = http.StatusCreated
		message = "Relação criada com sucesso"
	} else {
		status = http.StatusOK
		message = "Relação já existente"
	}

	utils.ToJson(w, status, utils.DefaultResponse{
		Data:    created,
		Message: message,
		Status:  status,
	})
}

func (h *UserCategoryRelationHandler) GetAllRelationsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetAllRelationsByUserID(ctx, id)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    relations,
		Message: "Relações recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryRelationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err1 := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	categoryID, err2 := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)

	if err1 != nil || err2 != nil {
		utils.ErrorResponse(w, fmt.Errorf("IDs inválidos"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, userID, categoryID); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserCategoryRelationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, userID); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
