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

const ref = "[AddressHandler] "

func (h *addressHandler) Enable(w http.ResponseWriter, r *http.Request) {
	h.handleToggle(
		w,
		r,
		ref+"[Enable] ",
		h.service.Enable,
	)
}

func (h *addressHandler) Disable(w http.ResponseWriter, r *http.Request) {
	h.handleToggle(
		w,
		r,
		ref+"[Disable] ",
		h.service.Disable,
	)
}

func (h *addressHandler) handleToggle(
	w http.ResponseWriter,
	r *http.Request,
	logRef string,
	action func(ctx context.Context, id int64) error,
) {
	ctx := r.Context()

	// valida método
	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, logRef+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(
			w,
			fmt.Errorf("método %s não permitido", r.Method),
			http.StatusMethodNotAllowed,
		)
		return
	}

	h.logger.Info(ctx, logRef+logger.LogUpdateInit, nil)

	// valida ID
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, logRef+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	// executa ação
	if err := action(ctx, id); err != nil {
		status := http.StatusInternalServerError
		clientErr := fmt.Errorf("erro interno")

		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			status = http.StatusBadRequest
			clientErr = fmt.Errorf("ID inválido")
			h.logger.Warn(ctx, logRef+logger.LogInvalidID, map[string]any{
				"address_id": id,
			})

		case errors.Is(err, errMsg.ErrNotFound):
			status = http.StatusNotFound
			clientErr = err
			h.logger.Warn(ctx, logRef+logger.LogNotFound, map[string]any{
				"address_id": id,
			})

		default:
			h.logger.Error(ctx, err, logRef+logger.LogUpdateError, map[string]any{
				"address_id": id,
			})
		}

		utils.ErrorResponse(w, clientErr, status)
		return
	}

	h.logger.Info(ctx, logRef+logger.LogUpdateSuccess, map[string]any{
		"address_id": id,
	})
	w.WriteHeader(http.StatusNoContent)
}
