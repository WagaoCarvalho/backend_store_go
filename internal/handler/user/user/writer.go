package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *User) Create(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Create] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData struct {
		User *dto.UserDTO `json:"user"`
	}

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userModel := dto.ToUserModel(*requestData.User)

	createdUser, err := h.service.Create(ctx, userModel)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"email": userModel.Email,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	})

	createdDTO := dto.ToUserDTO(createdUser)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *User) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Update] "
	ctx := r.Context()

	if r.Method != http.MethodPut {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var requestData struct {
		User *dto.UserDTO `json:"user"`
	}
	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	if requestData.User == nil {
		h.logger.Warn(ctx, ref+logger.LogMissingBodyData, nil)
		utils.ErrorResponse(w, fmt.Errorf("dados do usuário são obrigatórios"), http.StatusBadRequest)
		return
	}

	// Converte DTO para model
	userModel := dto.ToUserModel(*requestData.User)
	userModel.UID = id

	// Valida e atualiza via service
	if err := h.service.Update(ctx, userModel); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{"user_id": id})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return
		case errors.Is(err, errMsg.ErrInvalidData):
			h.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{"user_id": id, "erro": err.Error()})
			utils.ErrorResponse(w, err, http.StatusBadRequest)
			return
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"user_id": id})
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{"user_id": id})
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{"user_id": id})
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
	})
}

func (h *User) Delete(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Delete] "
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{})

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
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
				"user_id": id,
				"status":  status,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}
