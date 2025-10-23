package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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

func (h *Client) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var clientDTO dto.ClientDTO
	if err := utils.FromJSON(r.Body, &clientDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel := dto.ToClientModel(clientDTO)

	createdModel, err := h.service.Create(ctx, clientModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errMsg.ErrDBInvalidForeignKey, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"Cliente duplicado", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, errMsg.ErrDuplicate, http.StatusConflict)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	createdDTO := dto.ToClientDTO(createdModel)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"client_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Cliente criado com sucesso",
		Data:    createdDTO,
	})
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

func (h *Client) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - Update] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	uid, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var clientDTO dto.ClientDTO
	if err := utils.FromJSON(r.Body, &clientDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel := dto.ToClientModel(clientDTO)
	clientModel.ID = uid

	err = h.service.Update(ctx, clientModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidData),
			errors.Is(err, errMsg.ErrZeroID):
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
				"client_id": uid,
				"erro":      err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"Cliente duplicado", map[string]any{
				"client_id": uid,
				"erro":      err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	updatedDTO := dto.ToClientDTO(clientModel)

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"client_id": uid,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Cliente atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *Client) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - Delete] "
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
			"client_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"client_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *Client) Disable(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - Disable] "
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
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("cliente não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"client_id": uid,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Client) Enable(w http.ResponseWriter, r *http.Request) {
	const ref = "[clientHandler - Enable] "
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
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("cliente não encontrado"), http.StatusNotFound)
			return
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"client_id": uid,
			})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"client_id": uid,
	})
	w.WriteHeader(http.StatusNoContent)
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
