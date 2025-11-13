package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productCategoryRelationHandler) GetAllRelationsByProductID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - GetAllRelationsByProductID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "product_id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de usuário inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.productCategoryRelation.GetAllRelationsByProductID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// 🔹 Conversão para DTO
	var relationsDTO []dto.ProductCategoryRelationsDTO
	for _, rel := range relations {
		relationsDTO = append(relationsDTO, dto.ToProductCategoryRelationsDTO(rel))
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": id,
		"total":      len(relationsDTO),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    relationsDTO,
		Message: "Relações recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}
