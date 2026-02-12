package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *productHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductHandler - Create] "
	ctx := r.Context()

	// Validação de método HTTP
	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var productDTO dto.ProductDTO
	if err := utils.FromJSON(r.Body, &productDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	product := dto.ToProductModel(productDTO)

	createdProduct, err := h.service.Create(ctx, product)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+"chave estrangeira inválida", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("fornecedor inválido"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"produto duplicado", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("produto já existe"), http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrInvalidData):
			h.logger.Warn(ctx, ref+"dados inválidos", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
			utils.ErrorResponse(w, fmt.Errorf("erro ao criar produto"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"product_id": createdProduct.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Produto criado com sucesso",
		Data:    dto.ToProductDTO(createdProduct),
	})
}

func (h *productHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ProductHandler - Update] "
	ctx := r.Context()

	// Validação de método HTTP (aceita PUT ou PATCH)
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
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

	var productDTO dto.ProductDTO
	if err := utils.FromJSON(r.Body, &productDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	product := dto.ToProductModel(productDTO)
	product.ID = id

	// Opcional: garantir que o ID do DTO (se existir) seja ignorado
	productDTO.ID = nil

	err = h.service.Update(ctx, product)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidData),
			errors.Is(err, errMsg.ErrDBInvalidForeignKey),
			errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+"dados inválidos", map[string]any{
				"erro": err.Error(),
				"id":   id,
			})
			utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrVersionConflict),
			errors.Is(err, errMsg.ErrConflict):
			h.logger.Warn(ctx, ref+"conflito", map[string]any{
				"product_id": id,
				"erro":       err.Error(),
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de dados"), http.StatusConflict)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao atualizar produto"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": product.ID,
	})

	// Nota: product pode não ter dados atualizados do banco.
	// Se necessário, buscar produto atualizado ou retornar apenas confirmação
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Produto atualizado com sucesso",
		Data:    dto.ToProductDTO(product),
	})
}

func (h *productHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[productHandler - Delete] "
	ctx := r.Context()

	// Validação de método HTTP
	if r.Method != http.MethodDelete {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteInit, nil)

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
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("produto não encontrado"), http.StatusNotFound)
			return

		case errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
				"product_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro ao excluir produto"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"product_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
