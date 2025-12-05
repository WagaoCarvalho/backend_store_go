package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var userDTO dto.UserDTO
	if err := utils.FromJSON(r.Body, &userDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	userModel := dto.ToUserModel(userDTO)

	createdUser, err := h.service.Create(ctx, userModel)
	if err != nil {
		// CORREÇÃO: Verificar se o erro contém ErrInvalidData na string
		errStr := err.Error()
		if strings.Contains(errStr, errMsg.ErrInvalidData.Error()) {
			h.logger.Warn(ctx, ref+"Dados inválidos", map[string]any{
				"username": userModel.Username,
				"erro":     errStr,
			})

			// Extrai a mensagem após o prefixo do erro
			responseErr := errStr
			if strings.Contains(errStr, errMsg.ErrInvalidData.Error()+": ") {
				responseErr = strings.TrimPrefix(errStr, errMsg.ErrInvalidData.Error()+": ")
			}
			utils.ErrorResponse(w, fmt.Errorf(responseErr), http.StatusBadRequest)
			return
		}

		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"username": userModel.Username,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro interno do servidor"), http.StatusInternalServerError)
		return
	}

	// CORREÇÃO: Verificar se createdUser é nil antes de usar
	if createdUser == nil {
		h.logger.Error(ctx, errors.New("usuário criado é nulo"), ref+logger.LogCreateError, map[string]any{
			"username": userModel.Username,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro interno do servidor"), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
	})

	createdDTO := dto.ToUserDTO(createdUser)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Warn(ctx, ref+"ID inválido", map[string]any{
			"id_param": r.PathValue("id"),
			"erro":     err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var userDTO dto.UserDTO
	if err := utils.FromJSON(r.Body, &userDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	userModel := dto.ToUserModel(userDTO)
	userModel.UID = id

	h.logger.Info(ctx, ref+"Atualizando usuário", map[string]any{
		"user_id":  id,
		"username": userModel.Username,
		"version":  userModel.Version,
	})

	if err := h.service.Update(ctx, userModel); err != nil {
		errStr := err.Error()

		switch {
		case errors.Is(err, errMsg.ErrVersionConflict):
			h.logger.Warn(ctx, ref+"Conflito de versão", map[string]any{
				"user_id": id,
				"erro":    errStr,
			})
			utils.ErrorResponse(w, fmt.Errorf("versão desatualizada"), http.StatusConflict)
			return
		case errors.Is(err, errMsg.ErrInvalidData):
			h.logger.Warn(ctx, ref+"Dados inválidos", map[string]any{
				"user_id": id,
				"erro":    errStr,
			})
			// Extrai apenas a mensagem após o prefixo, se houver
			responseErr := errStr
			if strings.Contains(errStr, errMsg.ErrInvalidData.Error()+": ") {
				responseErr = strings.TrimPrefix(errStr, errMsg.ErrInvalidData.Error()+": ")
			}
			utils.ErrorResponse(w, fmt.Errorf(responseErr), http.StatusBadRequest)
			return
		case errors.Is(err, errMsg.ErrNotFound):
			h.logger.Warn(ctx, ref+"Usuário não encontrado", map[string]any{
				"user_id": id,
				"erro":    errStr,
			})
			utils.ErrorResponse(w, fmt.Errorf("usuário não encontrado"), http.StatusNotFound)
			return
		default:
			h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("erro interno do servidor"), http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": id,
	})

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
	})
}

func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
