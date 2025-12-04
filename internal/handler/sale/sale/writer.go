package handler

import (
	"errors"
	"fmt"
	"net/http"

	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *saleHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Create] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var saleDTO dtoSale.SaleDTO
	if err := utils.FromJSON(r.Body, &saleDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	saleModel := dtoSale.ToSaleModel(saleDTO)

	createdModel, err := h.service.Create(ctx, saleModel)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)

		switch {
		case errors.Is(err, errMsg.ErrInvalidData):
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			utils.ErrorResponse(w, err, http.StatusConflict)
			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	createdDTO := dtoSale.ToSaleDTO(createdModel)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{"sale_id": createdDTO.ID})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Venda criada com sucesso",
		Data:    createdDTO,
	})
}

func (h *saleHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Update] "
	ctx := r.Context()

	if r.Method != http.MethodPut {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+"Iniciando atualização da venda", nil)

	var saleDTO dtoSale.SaleDTO
	if err := utils.FromJSON(r.Body, &saleDTO); err != nil {
		h.logger.Warn(ctx, ref+"Erro ao parsear JSON", map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"id": id})
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	saleDTO.ID = &id
	saleModel := dtoSale.ToSaleModel(saleDTO)

	if err := h.service.Update(ctx, saleModel); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao atualizar venda", nil)

		switch {
		case errors.Is(err, errMsg.ErrInvalidData):
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrZeroVersion):
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Venda atualizada com sucesso", map[string]any{"sale_id": id})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Venda atualizada com sucesso",
		Data:    saleDTO,
	})
}

func (h *saleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SaleHandler - Delete] "
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{"method": r.Method})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+"Iniciando exclusão da venda", nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, errMsg.ErrZeroID, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao excluir venda", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{"sale_id": id})
	w.WriteHeader(http.StatusNoContent)
}
