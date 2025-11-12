package handler

import (
	"errors"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *clientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	clientDTO := dto.ToClientDTO(clientModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": clientDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Cliente encontrado",
		Data:    clientDTO,
	})
}

func (h *clientHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - GetByName] "
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

	clients, err := h.service.GetByName(ctx, name)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"name": name,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	clientDTOs := dto.ToClientDTOs(clients)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"name": name,
		"qtd":  len(clientDTOs),
	})

	message := "Clientes encontrados"
	if len(clientDTOs) == 0 {
		message = "Nenhum cliente encontrado"
	}

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    clientDTOs,
	})
}

func (h *clientHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - GetVersionByID] "
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
			"client_id": uid,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": uid,
		"version":   version,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do cliente recuperada com sucesso",
		Data: map[string]any{
			"client_id": uid,
			"version":   version,
		},
	})
}
