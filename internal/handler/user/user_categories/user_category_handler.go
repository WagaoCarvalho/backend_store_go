package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_categories"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
	"github.com/gorilla/mux"
)

type UserCategoryHandler struct {
	service service.UserCategoryService
	logger  *logger.LoggerAdapter
}

func NewUserCategoryHandler(service service.UserCategoryService, logger *logger.LoggerAdapter) *UserCategoryHandler {
	return &UserCategoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]interface{}{})

	var category *model.UserCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.Create(ctx, category)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]interface{}{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]interface{}{
		"category_id": createdCategory.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdCategory,
		Message: "Categoria criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *UserCategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryHandler - GetById] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]interface{}{})

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "categoria não encontrada" {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]interface{}{
				"id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]interface{}{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]interface{}{
		"id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    category,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]interface{}{})

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]interface{}{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]interface{}{
		"total": len(categories),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    categories,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]interface{}{})

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var category *model.UserCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	category.ID = uint(id)

	updatedCategory, err := h.service.Update(ctx, category)
	if err != nil {
		if errors.Is(err, repo.ErrCategoryNotFound) {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]interface{}{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]interface{}{
			"id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]interface{}{
		"id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    updatedCategory,
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ref := "[UserCategoryHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]interface{}{})

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]interface{}{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]interface{}{
		"id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
