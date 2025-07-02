package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
)

type UserCategoryHandler struct {
	service services.UserCategoryService
	logger  *logger.LoggerAdapter
}

func NewUserCategoryHandler(service services.UserCategoryService, logger *logger.LoggerAdapter) *UserCategoryHandler {
	return &UserCategoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryHandler] - Iniciando criação de categoria", nil)

	var category *models_user_categories.UserCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		h.logger.Warn(ctx, "[UserCategoryHandler] - Erro ao decodificar JSON da categoria", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.Create(ctx, category)
	if err != nil {
		h.logger.Error(ctx, err, "[UserCategoryHandler] - Erro ao criar categoria", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryHandler] - Categoria criada com sucesso", map[string]interface{}{
		"category_id": createdCategory.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdCategory,
		Message: "Categoria criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *UserCategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryHandler] - Iniciando busca por categoria pelo ID", nil)

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, "[UserCategoryHandler] - ID inválido recebido para busca", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "categoria não encontrada" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, "[UserCategoryHandler] - Categoria não encontrada", map[string]interface{}{
				"id": id,
			})
		} else {
			h.logger.Error(ctx, err, "[UserCategoryHandler] - Erro ao buscar categoria", map[string]interface{}{
				"id": id,
			})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, "[UserCategoryHandler] - Categoria recuperada com sucesso", map[string]interface{}{
		"id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    category,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryHandler] - Iniciando recuperação de todas as categorias", nil)

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Warn(ctx, "[UserCategoryHandler] - Erro ao recuperar categorias", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryHandler] - Categorias recuperadas com sucesso", map[string]interface{}{
		"total": len(categories),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    categories,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryHandler] - Iniciando atualização de categoria", nil)

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, "[UserCategoryHandler] - ID inválido para atualização", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var category *models_user_categories.UserCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		h.logger.Warn(ctx, "[UserCategoryHandler] - Falha ao decodificar JSON para atualização", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	category.ID = uint(id)

	updatedCategory, err := h.service.Update(ctx, category)
	if err != nil {
		if errors.Is(err, repo.ErrCategoryNotFound) {
			h.logger.Warn(ctx, "[UserCategoryHandler] - Categoria não encontrada para atualização", map[string]interface{}{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, "[UserCategoryHandler] - Erro ao atualizar categoria", map[string]interface{}{
			"id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryHandler] - Categoria atualizada com sucesso", map[string]interface{}{
		"id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    updatedCategory,
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "[UserCategoryHandler] - Iniciando exclusão de categoria", nil)

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, "[UserCategoryHandler] - ID inválido para exclusão", map[string]interface{}{
			"error": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, "[UserCategoryHandler] - Erro ao excluir categoria", map[string]interface{}{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "[UserCategoryHandler] - Categoria excluída com sucesso", map[string]interface{}{
		"id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
