package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *ProductCategory) GetByID(w http.ResponseWriter, r *http.Request) {
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

	productDTO := dto.ToProductCategoryDTO(category)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    productDTO,
		Message: "Categoria recuperada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *ProductCategory) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryHandler - GetAll] "
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

	productDTOs := dto.ToProductCategoryDTOs(categories)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    productDTOs,
		Message: "Categorias recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}
