package handlers

import (
	"net/http"

	supplier_models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type SupplierCategoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type supplierCategoryHandler struct {
	service supplier_categories.SupplierCategoryService
}

func NewSupplierCategoryHandler(service supplier_categories.SupplierCategoryService) SupplierCategoryHandler {
	return &supplierCategoryHandler{
		service: service,
	}
}

func (h *supplierCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category supplier_models.SupplierCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), &category)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	category.ID = id
	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Categoria de fornecedor criada com sucesso",
		Data:    category,
	})
}

func (h *supplierCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Categoria encontrada",
		Data:    category,
	})
}

func (h *supplierCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Categorias listadas com sucesso",
		Data:    categories,
	})
}

func (h *supplierCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var category supplier_models.SupplierCategory
	if err := utils.FromJson(r.Body, &category); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	category.ID = id

	if err := h.service.Update(r.Context(), &category); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Categoria atualizada com sucesso",
		Data:    category,
	})
}

func (h *supplierCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Categoria deletada com sucesso",
	})
}
