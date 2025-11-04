package handler

import (
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *Product) UpdateStock(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - UpdateStock] "
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
		Quantity int `json:"quantity"`
	}

	if err := utils.FromJSON(r.Body, &payload); err != nil {
		h.logger.Warn(ctx, ref+"erro ao decodificar payload", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateStock(ctx, uid, payload.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
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
		"quantity":   payload.Quantity,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Product) IncreaseStock(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - IncreaseStock] "
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
		Amount int `json:"stock_quantity"`
	}

	if err := utils.FromJSON(r.Body, &payload); err != nil {
		h.logger.Warn(ctx, ref+"erro ao decodificar payload", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.IncreaseStock(ctx, uid, payload.Amount)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
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
		"product_id":     uid,
		"stock_quantity": payload.Amount,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Product) DecreaseStock(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - DecreaseStock] "
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
		Amount int `json:"stock_quantity"`
	}

	if err := utils.FromJSON(r.Body, &payload); err != nil {
		h.logger.Warn(ctx, ref+"erro ao decodificar payload", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.DecreaseStock(ctx, uid, payload.Amount)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
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
		"product_id":     uid,
		"stock_quantity": payload.Amount,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Product) GetStock(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - GetStock] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+"iniciando", nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	stock, err := h.service.GetStock(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+"erro inesperado", map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	resp := map[string]any{
		"product_id":     uid,
		"stock_quantity": stock,
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, resp)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos listados com sucesso",
		Data:    resp,
	})
}
