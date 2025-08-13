package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
)

type SupplierCategoryRelationHandler struct {
	service services.SupplierCategoryRelationService
	logger  *logger.LoggerAdapter
}

func NewSupplierCategoryRelationHandler(service services.SupplierCategoryRelationService, logger *logger.LoggerAdapter) *SupplierCategoryRelationHandler {
	return &SupplierCategoryRelationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierCategoryRelationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryRelationHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var relation *models.SupplierCategoryRelations
	if err := utils.FromJson(r.Body, &relation); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	created, wasCreated, err := h.service.Create(ctx, relation.SupplierID, relation.CategoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidForeignKey) {
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": relation.SupplierID,
				"category_id": relation.CategoryID,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id": relation.SupplierID,
			"category_id": relation.CategoryID,
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
		"supplier_id": relation.SupplierID,
		"category_id": relation.CategoryID,
	})

	utils.ToJson(w, status, utils.DefaultResponse{
		Data:    created,
		Message: message,
		Status:  status,
	})
}

func (h *SupplierCategoryRelationHandler) GetBySupplierID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryRelationHandler - GetBySupplierID] "
	ctx := r.Context()

	idStr := mux.Vars(r)["supplier_id"]
	supplierID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": idStr})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetBySupplierId(ctx, supplierID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar relações", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações retornadas com sucesso", map[string]any{"supplier_id": supplierID})
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    relations,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryRelationHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryRelationHandler - GetByCategoryID] "
	ctx := r.Context()

	idStr := mux.Vars(r)["category_id"]
	categoryID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"category_id": idStr})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	relations, err := h.service.GetByCategoryId(ctx, categoryID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao buscar relações", map[string]any{"category_id": categoryID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Relações retornadas com sucesso", map[string]any{"category_id": categoryID})
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    relations,
		Message: "Relações encontradas",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryRelationHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryRelationHandler - DeleteByID] "
	ctx := r.Context()

	vars := mux.Vars(r)
	supplierID, err1 := strconv.ParseInt(vars["supplier_id"], 10, 64)
	categoryID, err2 := strconv.ParseInt(vars["category_id"], 10, 64)

	if err1 != nil || err2 != nil {
		h.logger.Warn(ctx, ref+"IDs inválidos", map[string]any{
			"supplier_id": vars["supplier_id"],
			"category_id": vars["category_id"],
		})
		utils.ErrorResponse(w, errors.New("IDs inválidos"), http.StatusBadRequest)
		return
	}

	err := h.service.DeleteById(ctx, supplierID, categoryID)
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
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Message: "Relação excluída com sucesso",
		Status:  http.StatusOK,
	})
}

func (h *SupplierCategoryRelationHandler) DeleteAllBySupplierID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierCategoryRelationHandler - DeleteAllBySupplierID] "
	ctx := r.Context()

	idStr := mux.Vars(r)["supplier_id"]
	supplierID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"supplier_id": idStr})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = h.service.DeleteAllBySupplierId(ctx, supplierID)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao excluir todas as relações", map[string]any{"supplier_id": supplierID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{"supplier_id": supplierID})
	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
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

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Data:    map[string]bool{"exists": hasRelation},
		Message: "Verificação concluída com sucesso",
		Status:  http.StatusOK,
	})
}
