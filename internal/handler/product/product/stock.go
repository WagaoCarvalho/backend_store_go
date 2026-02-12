package handler

import (
	"errors"
	"fmt"
	"net/http"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
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

	id, err := utils.GetIDParam(r, "id")
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

	err = h.service.UpdateStock(ctx, id, payload.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrInvalidQuantity):
			h.logger.Warn(ctx, ref+"quantidade inválida", map[string]any{
				"product_id": id,
				"quantity":   payload.Quantity,
			})
			utils.ErrorResponse(w, fmt.Errorf("quantidade inválida"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar estoque"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": id,
		"quantity":   payload.Quantity,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *productHandler) IncreaseStock(w http.ResponseWriter, r *http.Request) {
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

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var payload struct {
		Amount int `json:"amount"` // Campo mais semântico
	}

	if err := utils.FromJSON(r.Body, &payload); err != nil {
		h.logger.Warn(ctx, ref+"erro ao decodificar payload", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.IncreaseStock(ctx, id, payload.Amount)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrInvalidQuantity):
			h.logger.Warn(ctx, ref+"quantidade inválida", map[string]any{
				"product_id": id,
				"amount":     payload.Amount,
			})
			utils.ErrorResponse(w, fmt.Errorf("quantidade inválida"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao aumentar estoque"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": id,
		"amount":     payload.Amount,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *productHandler) DecreaseStock(w http.ResponseWriter, r *http.Request) {
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

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var payload struct {
		Amount int `json:"amount"`
	}

	if err := utils.FromJSON(r.Body, &payload); err != nil {
		h.logger.Warn(ctx, ref+"erro ao decodificar payload", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("payload inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.DecreaseStock(ctx, id, payload.Amount)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrInvalidQuantity):
			h.logger.Warn(ctx, ref+"quantidade inválida", map[string]any{
				"product_id": id,
				"amount":     payload.Amount,
			})
			utils.ErrorResponse(w, fmt.Errorf("quantidade inválida"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrInsufficientStock):
			h.logger.Warn(ctx, ref+"estoque insuficiente", map[string]any{
				"product_id": id,
				"amount":     payload.Amount,
			})
			utils.ErrorResponse(w, fmt.Errorf("estoque insuficiente"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao diminuir estoque"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": id,
		"amount":     payload.Amount,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *productHandler) GetStock(w http.ResponseWriter, r *http.Request) {
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

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	stock, err := h.service.GetStock(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+"erro inesperado", map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao buscar estoque"), http.StatusInternalServerError)
			return
		}
	}

	resp := map[string]any{
		"product_id":     id,
		"stock_quantity": stock,
	}

	h.logger.Info(ctx, ref+"sucesso", resp)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Estoque recuperado com sucesso",
		Data:    resp,
	})
}
