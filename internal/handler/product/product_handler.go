package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	service service.ProductService
	logger  *logger.LoggerAdapter
}

func NewProductHandler(service service.ProductService, logger *logger.LoggerAdapter) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - Create] "
	var product models.Product

	h.logger.Info(r.Context(), ref+logger.LogCreateInit, map[string]any{})

	if err := utils.FromJson(r.Body, &product); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.Create(r.Context(), &product)
	if err != nil {
		if errors.Is(err, repo.ErrInvalidForeignKey) {
			h.logger.Warn(r.Context(), ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		h.logger.Error(r.Context(), err, ref+logger.LogCreateError, map[string]any{})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogCreateSuccess, map[string]any{
		"product_id": createdProduct.ID,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Produto criado com sucesso",
		Data:    createdProduct,
	})
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ref := "[productHandler - GetAll] "

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"limit":  limit,
		"offset": offset,
	})

	products, err := h.service.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"limit":  limit,
			"offset": offset,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(products),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos listados com sucesso",
		Data:    products,
	})
}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - GetById] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := h.service.GetById(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": product.ID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto recuperado com sucesso",
		Data:    product,
	})
}

func (h *ProductHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - GetByName] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		h.logger.Warn(ctx, ref+logger.LogQueryError, map[string]any{
			"param": "name",
		})
		utils.ErrorResponse(w, errors.New("parâmetro 'name' é obrigatório"), http.StatusBadRequest)
		return
	}

	products, err := h.service.GetByName(ctx, name)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"name": name,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"result_count": len(products),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos recuperados com sucesso",
		Data:    products,
	})
}

func (h *ProductHandler) GetByManufacturer(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - GetByManufacturer] "
	ctx := r.Context()

	vars := mux.Vars(r)
	manufacturer := vars["manufacturer"]

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"manufacturer": manufacturer,
	})

	if manufacturer == "" {
		h.logger.Warn(ctx, ref+"manufacturer path param ausente", map[string]any{})
		utils.ErrorResponse(w, errors.New("parâmetro 'manufacturer' é obrigatório"), http.StatusBadRequest)
		return
	}

	products, err := h.service.GetByManufacturer(ctx, manufacturer)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "produtos não encontrados" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"manufacturer": manufacturer,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"manufacturer": manufacturer,
				"status":       status,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"manufacturer": manufacturer,
		"count":        len(products),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos encontrados",
		Data:    products,
	})
}

func (h *ProductHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - GetVersionByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(ctx, uid)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": uid,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": uid,
		"version":    version,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do produto recuperada com sucesso",
		Data: map[string]any{
			"product_id": uid,
			"version":    version,
		},
	})
}

func (h *ProductHandler) Disable(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - Disable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.Disable(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, repo.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Enable(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - Enable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.Enable(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, repo.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var input models.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Warn(ctx, ref+logger.LogMissingBodyData, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	input.ID = id

	updated, err := h.service.Update(ctx, &input)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": updated.ID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto atualizado com sucesso",
		Data:    updated,
	})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ref := "[productHandler - Delete] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"product_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"product_id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto deletado com sucesso",
	})
}
