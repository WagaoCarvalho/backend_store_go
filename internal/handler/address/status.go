package handler

import (
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address"
)

type Address struct {
	service service.Address
	logger  *logger.LogAdapter
}

func NewAddress(service service.Address, logger *logger.LogAdapter) *Address {
	return &Address{
		service: service,
		logger:  logger,
	}
}

func (h *Address) Enable(w http.ResponseWriter, r *http.Request) {
	const ref = "[AddressHandler - Enable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Enable(ctx, id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrZeroID) {
			status = http.StatusBadRequest
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"address_id": id,
			})
		} else if errors.Is(err, errMsg.ErrNotFound) {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"address_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"address_id": id,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"address_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *Address) Disable(w http.ResponseWriter, r *http.Request) {
	const ref = "[AddressHandler - Disable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Disable(ctx, id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrZeroID) {
			status = http.StatusBadRequest
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"address_id": id,
			})
		} else if errors.Is(err, errMsg.ErrNotFound) {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"address_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"address_id": id,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"address_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
