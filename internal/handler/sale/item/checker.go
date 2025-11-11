package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleItemHandler) ItemExists(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - ItemExists] "

	h.logger.Info(r.Context(), ref+logger.LogGetInit, nil)

	if r.Method != http.MethodGet {
		h.logger.Warn(r.Context(), ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		h.logger.Warn(r.Context(), ref+logger.LogInvalidID, map[string]any{
			"id":   id,
			"erro": errMsg.ErrZeroID.Error(),
		})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	exists, err := h.service.ItemExists(r.Context(), id)
	if err != nil {
		if errors.Is(err, errMsg.ErrDBInvalidForeignKey) {
			h.logger.Warn(r.Context(), ref+logger.LogForeignKeyViolation, map[string]any{
				"id":   id,
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, ref+logger.LogGetError, map[string]any{"id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	responseDTO := dto.ItemExistsResponseDTO{Exists: exists}

	h.logger.Info(r.Context(), ref+logger.LogGetSuccess, map[string]any{
		"id":     id,
		"exists": exists,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Verificação de existência realizada com sucesso",
		Data:    responseDTO,
	})
}
