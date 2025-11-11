package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - GetByID] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	item, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar item", map[string]any{"id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Converter para DTO
	itemDTO := dto.ToSaleItemDTO(item)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Item recuperado com sucesso",
		Data:    itemDTO,
	})
}

func (h *saleItemHandler) GetBySaleID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - GetBySaleID] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	saleID, err := utils.GetIDParam(r, "sale_id")
	if err != nil || saleID <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	limit, offset := utils.GetPaginationParams(r)

	items, err := h.service.GetBySaleID(ctx, saleID, limit, offset)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao listar itens por venda", map[string]any{"sale_id": saleID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Converter para DTO
	itemsDTO := dto.ToSaleItemDTOList(items)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Itens da venda recuperados com sucesso",
		Data:    itemsDTO,
	})
}

func (h *saleItemHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - GetByProductID] "
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	productID, err := utils.GetIDParam(r, "product_id")
	if err != nil || productID <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	limit, offset := utils.GetPaginationParams(r)

	items, err := h.service.GetByProductID(ctx, productID, limit, offset)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao listar itens por produto", map[string]any{"product_id": productID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Converter para DTO
	itemsDTO := dto.ToSaleItemDTOList(items)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Itens do produto recuperados com sucesso",
		Data:    itemsDTO,
	})
}
