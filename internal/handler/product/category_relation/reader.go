package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productCategoryRelationHandler) GetAllRelationsByProductID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductCategoryRelationHandler - GetAllRelationsByProductID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "product_id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de produto inválido"), http.StatusBadRequest)
		return
	}

	relations, err := h.productCategoryRelation.GetAllRelationsByProductID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Info(ctx, ref+"produto não encontrado", map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID de produto inválido"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao buscar relações do produto"), http.StatusInternalServerError)
			return
		}
	}

	// O serviço garante: relations nunca é nil (sempre slice vazio ou slice com valores)
	relationsDTO := make([]dto.ProductCategoryRelationDTO, 0, len(relations))
	for _, rel := range relations {
		relationsDTO = append(relationsDTO, dto.ToDTO(rel))
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": id,
		"count":      len(relationsDTO),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    relationsDTO,
		Message: "Relações recuperadas com sucesso",
		Status:  http.StatusOK,
	})
}
