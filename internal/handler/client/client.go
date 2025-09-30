package handler

import (
	"errors"
	"net/http"

	dtoClient "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client/client"
)

type ClientHandler struct {
	service service.ClientService
	logger  *logger.LogAdapter
}

func NewClientHandler(service service.ClientService, logger *logger.LogAdapter) *ClientHandler {
	return &ClientHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[ClientHandler - Create] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var clientDTO dtoClient.ClientDTO
	if err := utils.FromJSON(r.Body, &clientDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	clientModel := dtoClient.ToClientModel(clientDTO)

	createdModel, err := h.service.Create(ctx, clientModel)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrInvalidForeignKey):
			h.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return

		case errors.Is(err, errMsg.ErrDuplicate):
			h.logger.Warn(ctx, ref+"Cliente duplicado", map[string]any{
				"erro": err.Error(),
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return

		default:
			h.logger.Error(ctx, err, ref+logger.LogCreateError, nil)
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	createdDTO := dtoClient.ToClientDTO(createdModel)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"client_id": createdDTO.ID,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Cliente criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	clientDTO := dtoClient.ToClientDTO(clientModel)

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id": clientDTO.ID,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Cliente encontrado",
		Data:    clientDTO,
	})
}
