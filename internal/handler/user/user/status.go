package handler

import (
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *userHandler) Disable(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Disable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDisableInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	// Chama diretamente o service
	err = h.service.Disable(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogDisableError, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogDisableSuccess, map[string]any{
		"user_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) Enable(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Enable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogEnableInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	err = h.service.Enable(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogEnableError, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogEnableSuccess, map[string]any{
		"user_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
