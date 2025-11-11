package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *supplierCategoryRelationHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - GetBySupplierID] "
	ctx := r.Context()

	supplierID, err := utils.GetIDParam(r, "supplier_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetBySupplierID(ctx, supplierID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar relações", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações retornadas com sucesso", map[string]any{"supplier_id": supplierID})

	supplierDTO := dto.ToSupplierRelatinosDTOs(relations)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    supplierDTO,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *supplierCategoryRelationHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - GetByCategoryID] "
	ctx := r.Context()

	categoryID, err := utils.GetIDParam(r, "category_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"category_id": categoryID, "erro": err.Error()})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetByCategoryID(ctx, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar relações", map[string]any{"category_id": categoryID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações retornadas com sucesso", map[string]any{"category_id": categoryID})

	supplierDTO := dto.ToSupplierRelatinosDTOs(relations)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    supplierDTO,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *supplierCategoryRelationHandler) HasSupplierCategoryRelation(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - HasSupplierCategoryRelation] "

	supplierID, err := utils.GetIDParam(r, "supplier_id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"supplier_id inválido", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de fornecedor inválido: %w", err), http.StatusBadRequest)
		return
	}

	categoryID, err := utils.GetIDParam(r, "category_id")
	if err != nil {
		h.logger.Warn(r.Context(), ref+"category_id inválido", map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID de categoria inválido: %w", err), http.StatusBadRequest)
		return
	}

	h.logger.Info(r.Context(), ref+"iniciando verificação", map[string]any{
		"supplier_id": supplierID,
		"category_id": categoryID,
	})

	hasRelation, err := h.service.HasRelation(r.Context(), supplierID, categoryID)
	if err != nil {
		h.logger.Error(r.Context(), err, ref+"erro ao verificar relação", map[string]any{
			"supplier_id": supplierID,
			"category_id": categoryID,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao verificar relação: %w", err), http.StatusInternalServerError)
		return
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": hasRelation},
		Message: "Verificação concluída com sucesso",
		Status:  http.StatusOK,
	})
}
