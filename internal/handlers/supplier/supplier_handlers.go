package handlers

import (
	"fmt"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type SupplierHandler struct {
	service services.SupplierService
}

func NewSupplierHandler(service services.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

// CreateSupplier - POST /suppliers
func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var s models.Supplier
	if err := utils.FromJson(r.Body, &s); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), &s)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	s.ID = id
	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Fornecedor criado com sucesso",
		Data:    s,
	})
}

// GetSupplierByID - GET /suppliers/{id}
func (h *SupplierHandler) GetSupplierByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor encontrado",
		Data:    supplier,
	})
}

// GetAllSuppliers - GET /suppliers
func (h *SupplierHandler) GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	suppliers, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedores encontrados",
		Data:    suppliers,
	})
}

// UpdateSupplier - PUT /suppliers/{id}
func (h *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var s models.Supplier
	if err := utils.FromJson(r.Body, &s); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	s.ID = id
	if err := h.service.Update(r.Context(), &s); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor atualizado com sucesso",
		Data:    s,
	})
}

// DeleteSupplier - DELETE /suppliers/{id}
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor deletado com sucesso",
		Data:    nil,
	})
}
