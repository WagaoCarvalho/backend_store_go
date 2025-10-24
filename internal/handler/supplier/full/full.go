package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/full"
)

type Supplier struct {
	service service.SupplierFull
	logger  *logger.LogAdapter
}

func NewSupplierFull(service service.SupplierFull, logger *logger.LogAdapter) *Supplier {
	return &Supplier{
		service: service,
		logger:  logger,
	}
}

func (h *Supplier) CreateFull(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - CreateFull] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData dto.SupplierFullDTO

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Converte DTO para model antes de enviar para o serviço
	modelData := dto.ToSupplierFullModel(requestData)

	createdSupplierFull, err := h.service.CreateFull(ctx, modelData)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": modelData.Supplier.Name,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdSupplierFull.Supplier.ID,
		"name":        createdSupplierFull.Supplier.Name,
	})

	// Converte model de volta para DTO para resposta
	createdDTO := dto.ToSupplierFullDTO(createdSupplierFull)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Fornecedor criado com sucesso",
		Data:    createdDTO,
	})
}
