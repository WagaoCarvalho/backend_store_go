package handler

import (
	"fmt"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_full"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_full_services"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type SupplierHandler struct {
	service services.SupplierFullService
	logger  *logger.LoggerAdapter
}

func NewSupplierFullHandler(service services.SupplierFullService, logger *logger.LoggerAdapter) *SupplierHandler {
	return &SupplierHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierHandler) CreateFull(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - CreateFull] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData models.SupplierFull

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdSupplierFull, err := h.service.CreateFull(ctx, &requestData)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": requestData.Supplier.Name,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdSupplierFull.Supplier.ID,
		"name":        createdSupplierFull.Supplier.Name,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Fornecedor criado com sucesso",
		Data:    createdSupplierFull,
	})
}
