package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/contact_relation"
)

type SupplierContactRelation struct {
	service service.SupplierContactRelation
	logger  *logger.LogAdapter
}

func NewSupplierContactRelation(service service.SupplierContactRelation, logger *logger.LogAdapter) *SupplierContactRelation {
	return &SupplierContactRelation{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierContactRelation) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierContactRelationHandler - Create] "
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
		Relation *dto.ContactSupplierRelationDTO `json:"relation"`
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
	modelRelation := dto.ToContactSupplierRelationModel(*requestData.Relation)

	createdRelation, err := h.service.Create(ctx, modelRelation)
	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, errMsg.ErrZeroID):
			status = http.StatusBadRequest
		case errors.Is(err, errMsg.ErrRelationExists):
			status = http.StatusConflict
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			status = http.StatusBadRequest
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id": modelRelation.SupplierID,
			"contact_id":  modelRelation.ContactID,
		})
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdRelation.SupplierID,
		"contact_id":  createdRelation.ContactID,
	})

	createdDTO := dto.ToContactSupplierRelationDTO(createdRelation)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Relação criada com sucesso",
		Data:    createdDTO,
	})
}

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

func (h *SupplierContactRelation) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierContactRelationHandler - Delete] "
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

	if err := h.service.Delete(ctx, supplierID, contactID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar relação", map[string]any{
			"supplier_id": supplierID,
			"contact_id":  contactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relação deletada com sucesso", map[string]any{
		"supplier_id": supplierID,
		"contact_id":  contactID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    nil,
		Message: "Relação deletada com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *SupplierContactRelation) DeleteAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierContactRelationHandler - DeleteAll] "
	ctx := r.Context()

	supplierID, err := utils.GetIDParam(r, "supplier_id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAll(ctx, supplierID); err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao deletar relações", map[string]any{
			"supplier_id": supplierID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações deletadas com sucesso", map[string]any{
		"supplier_id": supplierID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    nil,
		Message: "Relações deletadas com sucesso",
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
