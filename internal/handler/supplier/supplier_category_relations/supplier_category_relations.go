package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier_category_relations"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_category_relations"
)

type SupplierCategoryRelationHandler struct {
	service service.SupplierCategoryRelationService
	logger  *logger.LogAdapter
}

func NewSupplierCategoryRelationHandler(service service.SupplierCategoryRelationService, logger *logger.LogAdapter) *SupplierCategoryRelationHandler {
	return &SupplierCategoryRelationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierCategoryRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData dto.SupplierCategoryRelationsDTO
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	modelRelation := dto.ToSupplierCategoryRelationsModel(requestData)

	created, wasCreated, err := h.service.Create(ctx, modelRelation.SupplierID, modelRelation.CategoryID)
	if err != nil {
		if errors.Is(err, errMsg.ErrDBInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": modelRelation.SupplierID,
				"category_id": modelRelation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id": modelRelation.SupplierID,
			"category_id": modelRelation.CategoryID,
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
		"category_id": modelRelation.CategoryID,
	})

	createdDTO := dto.ToSupplierCategoryRelationsDTO(created)

	utils.ToJSON(w, status, utils.DefaultResponse{
		Data:    createdDTO,
		Message: message,
		Status:  status,
	})
}

func (h *SupplierCategoryRelationHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
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

	supplierDTO := dto.ToSupplierRelatiosDTOs(relations)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    supplierDTO,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryRelationHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
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

	supplierDTO := dto.ToSupplierRelatiosDTOs(relations)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Data:    supplierDTO,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryRelationHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
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

	err := h.service.DeleteByID(ctx, supplierID, categoryID)
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

func (h *SupplierCategoryRelationHandler) DeleteAllBySupplierID(w http.ResponseWriter, r *http.Request) {
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

func (h *SupplierCategoryRelationHandler) HasSupplierCategoryRelation(w http.ResponseWriter, r *http.Request) {
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
