package handler

import (
	"errors"
	"fmt"
	"net/http"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
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
	const ref = "[addressHandler - GetByUserID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := h.service.GetByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+"usuário não encontrado", map[string]any{"user_id": id})
			utils.ErrorResponse(w, fmt.Errorf("usuário não encontrado"), http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"user_id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	addressDTOs := dtoAddress.ToAddressDTOs(addressModels)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": id,
		"count":   len(addressDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do usuário encontrados",
		Data:    addressDTOs,
	})
}

func (h *addressHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	const ref = "[addressHandler - GetByClientID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := h.service.GetByClientID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"client_id": id})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	addressDTOs := dtoAddress.ToAddressDTOs(addressModels)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": id,
		"count":     len(addressDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do cliente encontrados",
		Data:    addressDTOs,
	})
}

func (h *addressHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	const ref = "[addressHandler - GetBySupplierID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	addressModels, err := h.service.GetBySupplierID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"supplier_id": id})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	addressDTOs := dtoAddress.ToAddressDTOs(addressModels)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"count":       len(addressDTOs),
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereços do fornecedor encontrados",
		Data:    addressDTOs,
	})
}
