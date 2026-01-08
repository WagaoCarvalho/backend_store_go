package handler

import (
	"fmt"
	"net/http"
	"time"

	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByID] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"id": id})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	sale, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar venda", map[string]any{"id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	createdDTO := dtoSale.ToSaleDTO(sale)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Venda encontrada",
		Data:    createdDTO,
	})
}

func (h *saleHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByClientID] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	clientID, err := utils.GetIDParam(r, "client_id")
	if err != nil || clientID <= 0 {
		h.logger.Warn(ctx, ref+"clientID inválido", map[string]any{"client_id": clientID})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByClientID(ctx, clientID, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por clientID", map[string]any{"client_id": clientID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOs(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas do cliente recuperadas",
		Data:    salesDTO,
	})
}

func (h *saleHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByUserID] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.GetIDParam(r, "user_id")
	if err != nil || userID <= 0 {
		h.logger.Warn(ctx, ref+"userID inválido", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByUserID(ctx, userID, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por userID", map[string]any{"user_id": userID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOs(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas do usuário recuperadas",
		Data:    salesDTO,
	})
}

func (h *saleHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - GetByDateRange] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	startStr, errStart := utils.GetStringParam(r, "start")
	endStr, errEnd := utils.GetStringParam(r, "end")
	if errStart != nil || errEnd != nil {
		h.logger.Warn(ctx, ref+"datas ausentes ou inválidas", map[string]any{"start": startStr, "end": endStr})
		utils.ErrorResponse(w, errMsg.ErrInvalidData, http.StatusBadRequest)
		return
	}

	start, errStart := time.Parse(time.RFC3339, startStr)
	end, errEnd := time.Parse(time.RFC3339, endStr)
	if errStart != nil || errEnd != nil {
		h.logger.Warn(ctx, ref+"datas com formato inválido", map[string]any{"start": startStr, "end": endStr})
		utils.ErrorResponse(w, errMsg.ErrInvalidData, http.StatusBadRequest)
		return
	}

	limit, offset := utils.ParseLimitOffset(r)
	orderBy, orderDir := utils.ParseOrder(r)

	sales, err := h.service.GetByDateRange(ctx, start, end, limit, offset, orderBy, orderDir)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar vendas por período", map[string]any{"start": start, "end": end})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	salesDTO := dtoSale.ToSaleDTOs(sales)
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Vendas por período recuperadas",
		Data:    salesDTO,
	})
}
