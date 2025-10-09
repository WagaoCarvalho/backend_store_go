package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier_contact_relations"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_contact_relations"
)

type SupplierContactRelationHandler struct {
	service service.SupplierContactRelationServices
	logger  *logger.LogAdapter
}

func NewSupplierContactRelationHandler(service service.SupplierContactRelationServices, logger *logger.LogAdapter) *SupplierContactRelationHandler {
	return &SupplierContactRelationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierContactRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierContactRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData dto.ContactSupplierRelationDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToContactSupplierRelationModel(requestData)

	created, wasCreated, err := h.service.Create(ctx, modelRelation.SupplierID, modelRelation.ContactID)
	if err != nil {
		if errors.Is(err, errMsg.ErrInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": modelRelation.SupplierID,
				"contact_id":  modelRelation.ContactID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id": modelRelation.SupplierID,
			"contact_id":  modelRelation.ContactID,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	status := http.StatusOK
	message := "Relação já existente"
	logMsg := logger.LogAlreadyExists

	if wasCreated {
		status = http.StatusCreated
		message = "Relação criada com sucesso"
		logMsg = logger.LogCreateSuccess
	}

	h.logger.Info(ctx, ref+logMsg, map[string]any{
		"supplier_id": modelRelation.SupplierID,
		"contact_id":  modelRelation.ContactID,
	})

	createdDTO := dto.ToContactSupplierRelationDTO(created)

	utils.ToJSON(w, status, utils.DefaultResponse{
		Data:    createdDTO,
		Message: message,
		Status:  status,
	})
}

func (h *SupplierContactRelationHandler) GetAllBySupplierID(w http.ResponseWriter, r *http.Request) {
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

func (h *SupplierContactRelationHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *SupplierContactRelationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
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

func (h *SupplierContactRelationHandler) HasSupplierContactRelation(w http.ResponseWriter, r *http.Request) {
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
