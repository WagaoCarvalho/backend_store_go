package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *addressHandler) Enable(w http.ResponseWriter, r *http.Request) {
	h.handleToggle(
		w,
		r,
		"[AddressHandler - Enable] ",
		h.service.Enable,
	)
}

func (h *addressHandler) Disable(w http.ResponseWriter, r *http.Request) {
	h.handleToggle(
		w,
		r,
		"[AddressHandler - Disable] ",
		h.service.Disable,
	)
}

func (h *addressHandler) handleToggle(
	w http.ResponseWriter,
	r *http.Request,
	ref string,
	action func(ctx context.Context, id int64) error,
) {
	ctx := r.Context()

	// valida método
	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	// valida ID
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	// executa ação
	if err := action(ctx, id); err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			status = http.StatusBadRequest
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"address_id": id})

		case errors.Is(err, errMsg.ErrNotFound):
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"address_id": id})

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"address_id": id})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{"address_id": id})
	w.WriteHeader(http.StatusNoContent)
}
