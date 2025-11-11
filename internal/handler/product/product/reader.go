package handler

import (
	"errors"
	"net/http"
	"strconv"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
)

func (h *productHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[productHandler - GetAll] "

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

	productDTOs := dto.ToProductDTOs(products)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produtos listados com sucesso",
		Data:    productDTOs,
	})
}

func (h *productHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - GetByID] "
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

	product, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	productDTO := dto.ToProductDTO(product)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": product.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto recuperado com sucesso",
		Data:    productDTO,
	})
}

func (h *productHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - GetByName] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	name, err := utils.GetStringParam(r, "name")
	if err != nil || name == "" {
		h.logger.Warn(ctx, ref+logger.LogQueryError, map[string]any{
			"param": "name",
			"erro":  err,
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

	prouctDTOs := dto.ToProductDTOs(products)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"name":         name,
		"result_count": len(products),
	})

	message := "Produtos encontrados"
	if len(prouctDTOs) == 0 {
		message = "Nenhum produto encontrado"
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    prouctDTOs,
	})
}

func (h *productHandler) GetByManufacturer(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - GetByManufacturer] "
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
		if errors.Is(err, errMsg.ErrNotFound) {
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

	// converte para DTO
	productDTOs := dto.ToProductDTOs(products)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"manufacturer": manufacturer,
		"count":        len(productDTOs),
	})

	message := "Produtos encontrados"
	if len(productDTOs) == 0 {
		message = "Nenhum produto encontrado"
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    productDTOs,
	})
}

func (h *productHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - GetVersionByID] "
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

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do produto recuperada com sucesso",
		Data: map[string]any{
			"product_id": uid,
			"version":    version,
		},
	})
}
