package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_categories"
	"github.com/gorilla/mux"
)

type UserCategoryHandler struct {
	service service.UserCategoryService
	logger  *logger.LogAdapter
}

func NewUserCategoryHandler(service service.UserCategoryService, logger *logger.LogAdapter) *UserCategoryHandler {
	return &UserCategoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestDTO dto.UserCategoryDTO
	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	categoryModel := dto.ToUserCategoryModel(requestDTO)

	createdCategory, err := h.service.Create(ctx, categoryModel)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"category_id": createdCategory.ID,
	})

	createdDTO := dto.ToUserCategoryDTO(createdCategory)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdDTO,
		Message: "Categoria criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *UserCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - GetById] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "categoria não encontrada" {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"id": id,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    category,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(categories),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    categories,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var requestDTO dto.UserCategoryDTO

	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	modelCategory := dto.ToUserCategoryModel(requestDTO)
	modelCategory.ID = uint(id)

	updatedCategory, err := h.service.Update(ctx, modelCategory)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
		return
	}

	updatedDTO := dto.ToUserCategoryDTO(updatedCategory)

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"id": id,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    updatedDTO,
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *UserCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
