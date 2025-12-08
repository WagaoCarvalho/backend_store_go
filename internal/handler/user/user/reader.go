package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func (h *userHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetAll] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	users, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"quantidade": len(users),
	})

	userDTOs := dto.ToUserDTOs(users)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    userDTOs,
	})
}

func (h *userHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetByID] "
	ctx := r.Context()

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{})

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": id,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"user_id": id,
				"status":  status,
			})
		}

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": user.UID,
		"email":   user.Email,
	})

	userDTO := dto.ToUserDTO(user)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    userDTO,
	})
}

func (h *userHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetByEmail] "
	ctx := r.Context()

	email, err := utils.GetStringParam(r, "email")

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"email": email,
	})

	user, err := h.service.GetByEmail(ctx, email)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"email": email,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"email":  email,
				"status": status,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": user.UID,
		"email":   user.Email,
	})

	userDTO := dto.ToUserDTO(user)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    userDTO,
	})
}

func (h *userHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetByName] "
	ctx := r.Context()

	name, err := utils.GetStringParam(r, "username")

	h.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"name": name,
	})

	users, err := h.service.GetByName(ctx, name)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
			h.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"name": name,
			})
		} else {
			h.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
				"name":   name,
				"status": status,
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(users),
	})

	userDTOs := dto.ToUserDTOs(users)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    userDTOs,
	})
}
