package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user_categories"
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

	categories, err := h.service.GetCategories(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar categorias: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, fmt.Sprintf("erro ao encodar a resposta: %v", err), http.StatusInternalServerError)
	}
}

// GetCategoryById - Handler para buscar uma categoria por ID
func (h *UserCategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetCategoryById(ctx, id)
	if err != nil {
		if err.Error() == "categoria não encontrada" {
			http.Error(w, "categoria não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("erro ao buscar categoria: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		http.Error(w, fmt.Sprintf("erro ao encodar a resposta: %v", err), http.StatusInternalServerError)
	}
}

// CreateCategory - Handler para criar uma nova categoria
func (h *UserCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var category models.UserCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, fmt.Sprintf("erro ao decodificar o corpo da requisição: %v", err), http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.CreateCategory(ctx, category)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao criar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdCategory); err != nil {
		http.Error(w, fmt.Sprintf("erro ao encodar a resposta: %v", err), http.StatusInternalServerError)
	}
}

// UpdateCategory - Handler para atualizar uma categoria
func (h *UserCategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var category models.UserCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, fmt.Sprintf("erro ao decodificar o corpo da requisição: %v", err), http.StatusBadRequest)
		return
	}
	category.ID = uint(id) // Garantir que o ID seja atribuído

	updatedCategory, err := h.service.UpdateCategory(ctx, category)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedCategory); err != nil {
		http.Error(w, fmt.Sprintf("erro ao encodar a resposta: %v", err), http.StatusInternalServerError)
	}
}

// DeleteCategoryById - Handler para deletar uma categoria
func (h *UserCategoryHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCategoryById(ctx, id); err != nil {
		http.Error(w, fmt.Sprintf("erro ao deletar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
