package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier"
	"github.com/gorilla/mux"
)

type SupplierHandler struct {
	service service.SupplierService
	logger  *logger.LoggerAdapter
}

func NewSupplierHandler(service service.SupplierService, logger *logger.LoggerAdapter) *SupplierHandler {
	return &SupplierHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - Create] "
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
		Supplier *models.Supplier `json:"supplier"`
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdSupplier, err := h.service.Create(ctx, requestData.Supplier)
	if err != nil {
		// Evita acesso a ponteiro nulo
		logData := map[string]any{}
		if requestData.Supplier != nil {
			logData["name"] = requestData.Supplier.Name
			logData["cpf"] = requestData.Supplier.CPF
			logData["cnpj"] = requestData.Supplier.CNPJ
		}
		h.logger.Error(ctx, err, ref+logger.LogCreateError, logData)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdSupplier.ID,
		"name":        createdSupplier.Name,
		"cpf":         createdSupplier.CPF,
		"cnpj":        createdSupplier.CNPJ,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Fornecedor criado com sucesso",
		Data:    createdSupplier,
	})
}

func (h *SupplierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	suppliers, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar fornecedores: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(suppliers),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados",
		Data:    suppliers,
	})
}

func (h *SupplierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "fornecedor não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
				"status":      status,
			})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor encontrado",
		Data:    supplier,
	})
}

func (h *SupplierHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - GetVersionByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, err_msg.ErrNotFound) {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
				"status":      status,
			})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"version":     version,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do fornecedor obtida com sucesso",
		Data: map[string]int64{
			"version": version,
		},
	})
}

func (h *SupplierHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - GetByName] "
	ctx := r.Context()

	vars := mux.Vars(r)
	name := vars["name"]

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"name": name})

	suppliers, err := h.service.GetByName(ctx, name)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, err_msg.ErrNotFound) {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"name": name})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"name": name, "status": status})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(suppliers),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados com sucesso",
		Data:    suppliers,
	})
}

func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - Update] "
	ctx := r.Context()

	if r.Method != http.MethodPut {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	var requestData struct {
		Supplier *models.Supplier `json:"supplier"`
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
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

	requestData.Supplier.ID = id

	updatedSupplier, err := h.service.Update(ctx, requestData.Supplier)
	if err != nil {
		if errors.Is(err, err_msg.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": updatedSupplier.ID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor atualizado com sucesso",
		Data:    updatedSupplier,
	})
}

func (h *SupplierHandler) Disable(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - Disable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version <= 0 {
		h.logger.Warn(ctx, ref+"versão inválida", map[string]any{
			"erro": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "fornecedor não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	supplier.Status = false
	supplier.Version = payload.Version

	_, err = h.service.Update(ctx, supplier)
	if err != nil {
		if errors.Is(err, err_msg.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao desabilitar fornecedor: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *SupplierHandler) Enable(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - Enable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version <= 0 {
		h.logger.Warn(ctx, ref+"versão inválida", map[string]any{
			"erro": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "fornecedor não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"supplier_id": id,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	supplier.Status = true
	supplier.Version = payload.Version

	_, err = h.service.Update(ctx, supplier)
	if err != nil {
		if errors.Is(err, err_msg.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"supplier_id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao habilitar fornecedor: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ref := "[SupplierHandler - Delete] "
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
