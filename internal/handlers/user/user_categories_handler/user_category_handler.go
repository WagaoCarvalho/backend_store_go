package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
)

type UserCategoryHandler struct {
	service services.UserCategoryService
}

func NewUserCategoryHandler(service services.UserCategoryService) *UserCategoryHandler {
	return &UserCategoryHandler{service: service}
}

func (h *UserCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    categories,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "categoria não encontrada" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    category,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "categoria não encontrada" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]int{"version": version},
		Message: "Versão da categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var category *models_user_categories.UserCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.Create(ctx, category)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdCategory,
		Message: "Categoria criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *UserCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var category *models_user_categories.UserCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	category.ID = uint(id)

	updatedCategory, err := h.service.Update(ctx, category)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    updatedCategory,
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
