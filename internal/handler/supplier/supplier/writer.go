package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - Create] "
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
		Supplier *dto.SupplierDTO `json:"supplier"` // agora DTO
	}

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if requestData.Supplier == nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": "supplier não fornecido",
		})
		utils.ErrorResponse(w, fmt.Errorf("supplier não fornecido"), http.StatusBadRequest)
		return
	}

	// converte DTO para Model
	modelSupplier := dto.ToSupplierModel(*requestData.Supplier)

	createdSupplier, err := h.service.Create(ctx, modelSupplier)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": modelSupplier.Name,
			"cpf":  modelSupplier.CPF,
			"cnpj": modelSupplier.CNPJ,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdSupplier.ID,
		"name":        createdSupplier.Name,
		"cpf":         createdSupplier.CPF,
		"cnpj":        createdSupplier.CNPJ,
	})

	createdDTO := dto.ToSupplierDTO(createdSupplier)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Fornecedor criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - Update] "
	ctx := r.Context()

	if r.Method != http.MethodPut {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	// ID da URL
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	// Decodificar JSON
	var requestData struct {
		Supplier *dto.SupplierDTO `json:"supplier"`
	}

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	if requestData.Supplier == nil {
		h.logger.Warn(ctx, ref+logger.LogMissingBodyData, nil)
		utils.ErrorResponse(w, fmt.Errorf("dados do fornecedor são obrigatórios"), http.StatusBadRequest)
		return
	}

	// Setar ID vindo da URL
	if requestData.Supplier.ID == nil {
		requestData.Supplier.ID = new(int64)
	}
	*requestData.Supplier.ID = id

	supplierModel := dto.ToSupplierModel(*requestData.Supplier)

	// Chamar service
	err = h.service.Update(ctx, supplierModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidData),
			errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
				"supplier_id": id,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": id,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"Fornecedor duplicado", map[string]any{
				"supplier_id": id,
				"erro":        err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor atualizado com sucesso",
	})
}

func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[SupplierHandler - Delete] "
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "fornecedor não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
				"supplier_id": id,
				"status":      status,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
