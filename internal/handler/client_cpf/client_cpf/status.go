package handler

import (
	"context"
	"errors"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *clientCpfHandler) toggleStatus(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, id int64) error,
	ref string,
) {
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, errMsg.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, nil)
		utils.ErrorResponse(w, errMsg.ErrInvalidID, http.StatusBadRequest)
		return
	}

	if err := action(ctx, id); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"client_id": id})
			utils.ErrorResponse(w, errMsg.ErrNotFound, http.StatusNotFound)
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"version conflict", map[string]any{"client_id": id})
			utils.ErrorResponse(w, errMsg.ErrVersionConflict, http.StatusConflict)
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"client_id": id})
			utils.ErrorResponse(w, errMsg.ErrUpdate, http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{"client_id": id})
	w.WriteHeader(http.StatusNoContent)
}

func (h *clientCpfHandler) Disable(w http.ResponseWriter, r *http.Request) {
	h.toggleStatus(w, r, h.service.Disable, "[ClientHandler - Disable] ")
}

func (h *clientCpfHandler) Enable(w http.ResponseWriter, r *http.Request) {
	h.toggleStatus(w, r, h.service.Enable, "[ClientHandler - Enable] ")
}
