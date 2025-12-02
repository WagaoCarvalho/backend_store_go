package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - Create] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var itemDTO dto.SaleItemDTO
	if err := utils.FromJSON(r.Body, &itemDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	itemModel := dto.ToSaleItemModel(itemDTO)
	createdItem, err := h.service.Create(ctx, itemModel)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)

		status := http.StatusInternalServerError
		if errors.Is(err, errMsg.ErrInvalidData) || errors.Is(err, errMsg.ErrDBInvalidForeignKey) {
			status = http.StatusBadRequest
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	createdDTO := dto.ToSaleItemDTO(createdItem)
	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{"sale_item_id": createdDTO.ID})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Item de venda criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *saleItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - Update] "
	ctx := r.Context()

	if r.Method != http.MethodPut {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var itemDTO dto.SaleItemDTO
	if err := utils.FromJSON(r.Body, &itemDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	item := dto.ToSaleItemModel(itemDTO)
	item.ID = id

	err = h.service.Update(ctx, item)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidData),
			errors.Is(err, errMsg.ErrDBInvalidForeignKey),
			errors.Is(err, errMsg.ErrZeroID):
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroVersion):
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"sale_item_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"sale_item_id": item.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Item de venda atualizado com sucesso",
		Data:    dto.ToSaleItemDTO(item),
	})
}

func (h *saleItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - Delete] "
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar item de venda", map[string]any{"id": id})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *saleItemHandler) DeleteBySaleID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleItemHandler - DeleteBySaleID] "
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	saleID, err := utils.GetIDParam(r, "sale_id")
	if err != nil || saleID <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBySaleID(ctx, saleID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar itens da venda", map[string]any{"sale_id": saleID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
