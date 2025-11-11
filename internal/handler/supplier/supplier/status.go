package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *supplierHandler) Disable(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - Disable] "
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

	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version <= 0 {
		h.logger.Warn(ctx, ref+"versão inválida", map[string]any{
			"erro": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	supplier.Status = false
	supplier.Version = payload.Version

	err = h.service.Update(ctx, supplier)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao desabilitar fornecedor: %w", err), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *supplierHandler) Enable(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - Enable] "
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

	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version <= 0 {
		h.logger.Warn(ctx, ref+"versão inválida", map[string]any{
			"erro": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("fornecedor não encontrado"), http.StatusNotFound)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar fornecedor: %w", err), http.StatusInternalServerError)
		return
	}

	supplier.Status = true
	supplier.Version = payload.Version

	if err := h.service.Update(ctx, supplier); err != nil {
		if errors.Is(err, errMsg.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao habilitar fornecedor: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
