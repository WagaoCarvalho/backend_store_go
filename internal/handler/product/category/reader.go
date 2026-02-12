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

func (h *productCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - GetById] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(ctx, int64(id))
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("categoria não encontrada"), http.StatusNotFound)
		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+"ID zero", map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		default:
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao buscar categoria"), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"id": id,
	})

	productDTO := dto.ToProductCategoryDTO(category)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    productDTO,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *productCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar categorias"), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(categories),
	})

	productDTOs := dto.ToProductCategoryDTOs(categories)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    productDTOs,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}
