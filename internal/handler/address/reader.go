package handler

import (
	"context"
	"errors"
	"net/http"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
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
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
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

func (h *addressHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	h.handleGetAddresses(
		w, r,
		"[addressHandler - GetByUserID] ",
		"user_id",
		h.service.GetByUserID,
		"Endereços do usuário encontrados",
	)
}

func (h *addressHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	h.handleGetAddresses(
		w, r,
		"[addressHandler - GetByClientID] ",
		"client_id",
		h.service.GetByClientID,
		"Endereços do cliente encontrados",
	)
}

func (h *addressHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	h.handleGetAddresses(
		w, r,
		"[addressHandler - GetBySupplierID] ",
		"supplier_id",
		h.service.GetBySupplierID,
		"Endereços do fornecedor encontrados",
	)
}

func (h *addressHandler) handleGetAddresses(
	w http.ResponseWriter,
	r *http.Request,
	ref string,
	idLabel string,
	getFn func(ctx context.Context, id int64) ([]*models.Address, error),
	successMsg string,
) {
	ctx := r.Context()
	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := getFn(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogGetError, map[string]any{idLabel: id})
			utils.ErrorResponse(w, errMsg.ErrNotFound, http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{idLabel: id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	dto := dtoAddress.ToAddressDTOs(addressModels)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		idLabel: id,
		"count": len(dto),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: successMsg,
		Data:    dto,
	})
}
