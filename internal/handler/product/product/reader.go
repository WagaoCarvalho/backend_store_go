package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, errors.New("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, errors.New("erro ao buscar produto"), http.StatusInternalServerError)
			return
		}
	}

	productDTO := dto.ToProductDTO(product)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": product.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto recuperado com sucesso",
		Data:    productDTO,
	})
}
