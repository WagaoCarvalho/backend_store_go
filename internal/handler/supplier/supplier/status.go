package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

const ref = "[SupplierHandler] "

func (h *supplierHandler) Enable(w http.ResponseWriter, r *http.Request) {
	h.handleToggle(
		w,
		r,
		ref+"[Enable] ",
		h.service.Enable,
		true, // status para habilitar
	)
}

func (h *supplierHandler) Disable(w http.ResponseWriter, r *http.Request) {
	h.handleToggle(
		w,
		r,
		ref+"[Disable] ",
		h.service.Disable,
		false, // status para desabilitar
	)
}

func (h *supplierHandler) handleToggle(
	w http.ResponseWriter,
	r *http.Request,
	logRef string,
	action func(ctx context.Context, id int64) error,
	desiredStatus bool,
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

	// valida versão do payload
	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Warn(ctx, logRef+"invalid payload", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	if payload.Version <= 0 {
		h.logger.Warn(ctx, logRef+"invalid version", map[string]any{
			"version": payload.Version,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
		return
	}

	// busca o fornecedor atual para verificar versão
	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		clientErr := fmt.Errorf("erro interno")

		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			status = http.StatusNotFound
			clientErr = err
			h.logger.Warn(ctx, logRef+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})

		default:
			h.logger.Error(ctx, err, logRef+logger.LogGetError, map[string]any{
				"supplier_id": id,
			})
		}

		utils.ErrorResponse(w, clientErr, status)
		return
	}

	// verifica conflito de versão
	if supplier.Version != payload.Version {
		h.logger.Warn(ctx, logRef+"version conflict", map[string]any{
			"supplier_id":     id,
			"current_version": supplier.Version,
			"request_version": payload.Version,
		})
		utils.ErrorResponse(
			w,
			fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"),
			http.StatusConflict,
		)
		return
	}

	// executa ação (enable/disable)
	if err := action(ctx, id); err != nil {
		status := http.StatusInternalServerError
		clientErr := fmt.Errorf("erro interno")

		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			status = http.StatusNotFound
			clientErr = err
			h.logger.Warn(ctx, logRef+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})

		default:
			h.logger.Error(ctx, err, logRef+logger.LogUpdateError, map[string]any{
				"supplier_id": id,
			})
		}

		utils.ErrorResponse(w, clientErr, status)
		return
	}

	h.logger.Info(ctx, logRef+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
		"status":      desiredStatus,
	})

	w.WriteHeader(http.StatusNoContent)
}
