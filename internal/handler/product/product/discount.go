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

func (h *productHandler) EnableDiscount(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - EnableDiscount] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.EnableDiscount(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *productHandler) DisableDiscount(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - DisableDiscount] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.DisableDiscount(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *productHandler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - ApplyDiscount] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var payload struct {
		Percent float64 `json:"percent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Warn(ctx, ref+"payload inválido", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	product, err := h.service.ApplyDiscount(ctx, uid, payload.Percent)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
		"percent":    payload.Percent,
	})

	utils.ToJSON(w, http.StatusOK, product)
}
