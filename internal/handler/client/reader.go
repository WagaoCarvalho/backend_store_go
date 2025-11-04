package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client/client"
)

type Client struct {
	service service.Client
	logger  *logger.LogAdapter
}

func NewClient(service service.Client, logger *logger.LogAdapter) *Client {
	return &Client{
		service: service,
		logger:  logger,
	}
}

func (h *Client) GetByID(w http.ResponseWriter, r *http.Request) {
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

func (h *Client) GetByName(w http.ResponseWriter, r *http.Request) {
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

func (h *Client) GetVersionByID(w http.ResponseWriter, r *http.Request) {
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

func (h *Client) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const ref = "[clientHandler - GetAll] "

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

	clients, err := h.service.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"limit":  limit,
			"offset": offset,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(clients),
	})

	clientDTOs := dto.ToClientDTOs(clients)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Clientes listados com sucesso",
		Data:    clientDTOs,
	})
}

func (h *Client) ClientExists(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - ClientExists] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, nil)

	clientID, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	exists, err := h.service.ClientExists(ctx, clientID)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogNotFound, map[string]any{"client_id": clientID})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Verificação concluída", map[string]any{
		"client_id": clientID,
		"exists":    exists,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Verificação concluída com sucesso",
		Data: map[string]any{
			"client_id": clientID,
			"exists":    exists,
		},
	})
}
