package handler

import (
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/contact_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *SupplierContactRelation) GetAllBySupplierID(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierContactRelationHandler - GetAllBySupplierID] "
	ctx := r.Context()

	supplierID, err := utils.GetIDParam(r, "supplier_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetAllRelationsBySupplierID(ctx, supplierID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar relações", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações retornadas com sucesso", map[string]any{"supplier_id": supplierID})

	relationsDTO := dto.ToSupplierContactRelationsDTOs(relations)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    relationsDTO,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *SupplierContactRelation) HasSupplierContactRelation(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierContactRelationHandler - HasRelation] "
	ctx := r.Context()

	supplierID, err := utils.GetIDParam(r, "supplier_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	contactID, err := utils.GetIDParam(r, "contact_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"contact_id": contactID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	exists, err := h.service.HasSupplierContactRelation(ctx, supplierID, contactID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao verificar relação", map[string]any{
			"supplier_id": supplierID,
			"contact_id":  contactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Verificação concluída", map[string]any{
		"supplier_id": supplierID,
		"contact_id":  contactID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": exists},
		Message: "Verificação concluída",
		Status:  http.StatusOK,
	})
}
