package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *supplierCategoryRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - Create] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData struct {
		Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
	}

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if requestData.Relation == nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": "relation não fornecida",
		})
		utils.ErrorResponse(w, fmt.Errorf("relation não fornecida"), http.StatusBadRequest)
		return
	}

	// converte DTO para Model
	modelRelation := dto.ToSupplierCategoryRelationsModel(*requestData.Relation)

	createdRelation, err := h.service.Create(ctx, modelRelation)
	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			status = http.StatusBadRequest
		case errors.Is(err, errMsg.ErrRelationExists):
			status = http.StatusConflict
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey): // <--- adicionar
			status = http.StatusBadRequest
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id": modelRelation.SupplierID,
			"category_id": modelRelation.CategoryID,
		})
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdRelation.SupplierID,
		"category_id": createdRelation.CategoryID,
	})

	createdDTO := dto.ToSupplierCategoryRelationsDTO(createdRelation)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Relação criada com sucesso",
		Data:    createdDTO,
	})
}

func (h *supplierCategoryRelationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - DeleteByID] "
	ctx := r.Context()

	supplierID, err1 := utils.GetIDParam(r, "supplier_id")
	categoryID, err2 := utils.GetIDParam(r, "category_id")

	if err1 != nil || err2 != nil {
		h.logger.Warn(ctx, ref+"IDs inválidos", map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, errors.New("IDs inválidos"), http.StatusBadRequest)
		return
	}

	err := h.service.Delete(ctx, supplierID, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao excluir relação", map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Message: "Relação excluída com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *supplierCategoryRelationHandler) DeleteAllBySupplierID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - DeleteAllBySupplierID] "
	ctx := r.Context()

	supplierID, err := utils.GetIDParam(r, "supplier_id")

	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = h.service.DeleteAllBySupplierID(ctx, supplierID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao excluir todas as relações", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{"supplier_id": supplierID})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Message: "Relações excluídas com sucesso",
		Status:  http.StatusOK,
	})
}
