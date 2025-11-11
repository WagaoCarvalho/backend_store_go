package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *userCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
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

func (h *userCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.Update(ctx, modelCategory); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+"ID inválido", map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrInvalidData):
			h.logger.Warn(ctx, ref+"dados inválidos", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"id": id,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *userCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserCategoryHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
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
