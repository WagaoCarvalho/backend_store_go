package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/products"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.service.GetProducts(ctx)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar produtos: %w", err), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Data:   products,
		Status: http.StatusOK,
	}

	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductById(r.Context(), id)
	if err != nil {
		if err.Error() == "produto não encontrado" {
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
		} else {
			utils.ErrorResponse(w, fmt.Errorf("erro ao buscar produto"), http.StatusInternalServerError)
		}
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto encontrado",
		Data:    product,
	}
	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), product)
	if err != nil {
		if strings.Contains(err.Error(), "validação falhou") {
			utils.ErrorResponse(w, err, http.StatusBadRequest)
		} else {
			utils.ErrorResponse(w, fmt.Errorf("erro ao criar produto"), http.StatusInternalServerError)
		}
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Produto criado com sucesso",
		Data:    createdProduct,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	product.ID = int(id)

	updatedProduct, err := h.service.UpdateProduct(r.Context(), product)
	if err != nil {
		if err.Error() == "produto não encontrado" {
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
		} else {
			utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar produto"), http.StatusInternalServerError)
		}
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto atualizado com sucesso",
		Data:    updatedProduct,
	}
	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProductById(r.Context(), id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "produto não encontrado") {
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
		} else {
			utils.ErrorResponse(w, fmt.Errorf("erro ao deletar produto"), http.StatusInternalServerError)
		}
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto deletado com sucesso",
		Data:    nil,
	}
	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) GetProductsByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		utils.ErrorResponse(w, fmt.Errorf("parâmetro 'name' é obrigatório"), http.StatusBadRequest)
		return
	}

	products, err := h.service.GetProductsByName(r.Context(), name)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar produtos por nome"), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos encontrados",
		Data:    products,
	}
	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) GetProductsByCostPriceRange(w http.ResponseWriter, r *http.Request) {
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")

	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("valor mínimo inválido"), http.StatusBadRequest)
		return
	}

	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("valor máximo inválido"), http.StatusBadRequest)
		return
	}

	products, err := h.service.GetProductsByCostPriceRange(r.Context(), min, max)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar produtos por faixa de preço de custo"), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos encontrados por faixa de preço de custo",
		Data:    products,
	}
	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) GetProductsBySalePriceRange(w http.ResponseWriter, r *http.Request) {
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")

	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("valor mínimo inválido"), http.StatusBadRequest)
		return
	}

	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("valor máximo inválido"), http.StatusBadRequest)
		return
	}

	products, err := h.service.GetProductsBySalePriceRange(r.Context(), min, max)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar produtos por faixa de preço de venda"), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos encontrados por faixa de preço de venda",
		Data:    products,
	}
	utils.ToJson(w, http.StatusOK, response)
}

func (h *ProductHandler) GetProductsLowInStock(w http.ResponseWriter, r *http.Request) {
	thresholdStr := r.URL.Query().Get("threshold")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("valor de threshold inválido"), http.StatusBadRequest)
		return
	}

	products, err := h.service.GetProductsLowInStock(r.Context(), threshold)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar produtos com estoque baixo"), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos com estoque baixo encontrados",
		Data:    products,
	}
	utils.ToJson(w, http.StatusOK, response)
}
