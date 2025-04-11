package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
)

type UserCategoryHandler struct {
	service services.UserCategoryService
}

func NewUserCategoryHandler(service services.UserCategoryService) *UserCategoryHandler {
	return &UserCategoryHandler{service: service}
}

// GetCategories - Handler para buscar todas as categorias
func (h *UserCategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
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

// GetCategoryById - Handler para buscar uma categoria por ID
func (h *UserCategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetById(ctx, id)
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

// CreateCategory - Handler para criar uma nova categoria
func (h *UserCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var category models.UserCategory
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

// UpdateCategory - Handler para atualizar uma categoria
func (h *UserCategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var category models.UserCategory
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

// DeleteCategoryById - Handler para deletar uma categoria
func (h *UserCategoryHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
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
