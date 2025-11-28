package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestDTO dto.ProductCategoryDTO
	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	categoryModel := dto.ToProductCategoryModel(requestDTO)

	createdCategory, err := h.service.Create(ctx, categoryModel)
	if err != nil {
		if errors.Is(err, errMsg.ErrAlreadyExists) {
			h.logger.Warn(ctx, ref+"categoria já existe", nil)
			utils.ErrorResponse(w, fmt.Errorf("categoria já existe"), http.StatusConflict)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	createdDTO := dto.ToProductCategoryDTO(createdCategory)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Data:    createdDTO,
		Message: "Categoria criada com sucesso",
		Status:  http.StatusCreated,
	})
}

func (h *productCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var requestDTO dto.ProductCategoryDTO
	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	category := dto.ToProductCategoryModel(requestDTO)
	category.ID = uint(id)

	err = h.service.Update(ctx, category)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"id": id})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"id": id})
			utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar categoria: %v", err), http.StatusInternalServerError)
		}
		return
	}

	updatedCategory, err := h.service.GetByID(ctx, int64(id))
	if err != nil {
		h.logger.Error(ctx, err, ref+"erro ao buscar categoria atualizada", map[string]any{"id": id})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar categoria atualizada"), http.StatusInternalServerError)
		return
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    dto.ToProductCategoryDTO(updatedCategory),
		Message: "Categoria atualizada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *productCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - Delete] "
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
