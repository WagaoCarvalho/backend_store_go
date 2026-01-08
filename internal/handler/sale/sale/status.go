package handler

import (
	"fmt"
	"net/http"

	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByStatus] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	status, err := utils.GetStringParam(r, "status")
	if err != nil {
		h.logger.Warn(ctx, ref+"status inválido", map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, errMsg.ErrInvalidData, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByStatus(ctx, status, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por status", map[string]any{"status": status})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOs(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas por status recuperadas",
		Data:    salesDTO,
	})
}

func (h *saleHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Cancel] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.Cancel(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao cancelar venda", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *saleHandler) Complete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Complete] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.Complete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao completar venda", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *saleHandler) Returned(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Returned] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.Returned(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao marcar venda como devolvida", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *saleHandler) Activate(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Activate] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.Activate(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao reativar venda", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
