package handler

import (
	"errors"
	"net/http"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *addressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[addressHandler - GetByID] "
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

	addressModel, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogGetError, map[string]any{
				"address_id": id,
			})
			utils.ErrorResponse(w, errMsg.ErrNotFound, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	addressDTO := dtoAddress.ToAddressDTO(addressModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"address_id": addressDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço encontrado",
		Data:    addressDTO,
	})
}
