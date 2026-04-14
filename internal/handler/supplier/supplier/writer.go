package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *supplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[SupplierHandler] "

	h.logger.Info(ctx, ref+"[Create] "+logger.LogCreateInit, nil)

	var supplierDTO dto.SupplierDTO
	if err := utils.FromJSON(r.Body, &supplierDTO); err != nil {
		h.logger.Warn(ctx, ref+"[Create] "+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("JSON inválido"), http.StatusBadRequest)
		return
	}

	supplierModel := dto.ToSupplierModel(supplierDTO)

	// Validação dos dados (se o model tiver método Validate)
	if err := supplierModel.Validate(); err != nil {
		h.logger.Warn(ctx, ref+"[Create] "+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("dados inválidos"), http.StatusUnprocessableEntity)
		return
	}

	createdModel, err := h.service.Create(ctx, supplierModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"[Create] Fornecedor duplicado", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errors.New("fornecedor já existente"), http.StatusConflict)

		default:
			h.logger.Error(ctx, err, ref+"[Create] "+logger.LogCreateError, nil)
			utils.ErrorResponse(w, errors.New("erro interno"), http.StatusInternalServerError)
		}
		return
	}

	createdDTO := dto.ToSupplierDTO(createdModel)

	h.logger.Info(ctx, ref+"[Create] "+logger.LogCreateSuccess, map[string]any{
		"supplier_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Fornecedor criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *supplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"[Update] "+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)
		return
	}

	var supplierDTO dto.SupplierDTO // ← Certifique-se que o pacote dto está importado
	if err := utils.FromJSON(r.Body, &supplierDTO); err != nil {
		h.logger.Warn(ctx, ref+"[Update] "+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("JSON inválido"), http.StatusBadRequest)
		return
	}

	// CORRIGIDO: O version vem do DTO, não precisa ser definido manualmente
	supplierDTO.ID = &id
	supplierModel := dto.ToSupplierModel(supplierDTO)

	// Validação dos dados
	if err := supplierModel.Validate(); err != nil {
		h.logger.Warn(ctx, ref+"[Update] "+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
			"id":   id,
		})
		utils.ErrorResponse(w, errors.New("dados inválidos"), http.StatusUnprocessableEntity)
		return
	}

	h.logger.Info(ctx, ref+"[Update] "+logger.LogUpdateInit, map[string]any{
		"supplier_id": id,
		"version":     supplierModel.Version, // ← Log para depuração
	})

	if err := h.service.Update(ctx, supplierModel); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+"[Update] "+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errors.New("fornecedor não encontrado"), http.StatusNotFound)

		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"[Update] Conflito de versão", map[string]any{
				"supplier_id": id,
				"version":     supplierModel.Version,
			})
			utils.ErrorResponse(w, errors.New("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"[Update] Fornecedor duplicado", map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errors.New("fornecedor já existente"), http.StatusConflict)

		default:
			h.logger.Error(ctx, err, ref+"[Update] "+logger.LogUpdateError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errors.New("erro interno"), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(ctx, ref+"[Update] "+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": id,
		"new_version": supplierModel.Version, // ← Log da nova versão
	})

	updatedDTO := dto.ToSupplierDTO(supplierModel)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Fornecedor atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *supplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"[Delete] "+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, errors.New("ID inválido"), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, ref+"[Delete] "+logger.LogDeleteInit, map[string]any{
		"supplier_id": id,
		"path":        r.URL.Path,
	})

	if err := h.service.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+"[Delete] "+logger.LogNotFound, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errors.New("fornecedor não encontrado"), http.StatusNotFound)

		default:
			h.logger.Error(ctx, err, ref+"[Delete] "+logger.LogDeleteError, map[string]any{
				"supplier_id": id,
			})
			utils.ErrorResponse(w, errors.New("erro interno"), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(ctx, ref+"[Delete] "+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
