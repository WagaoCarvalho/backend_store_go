package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/users"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service services.UserService
	logger  *logger.LoggerAdapter
}

func NewUserHandler(service services.UserService, logger *logger.LoggerAdapter) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - Create] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{})

	var requestData struct {
		User *models.User `json:"user"`
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.Create(ctx, requestData.User)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"email": requestData.User.Email,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUser.UID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdUser,
	})
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - GetAll] "
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

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    users,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - GetByID] "
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

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	})
}

func (h *UserHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - GetVersionByID] "
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

	version, err := h.service.GetVersionByID(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrUserNotFound) {
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
		"user_id": id,
		"version": version,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do usuário obtida com sucesso",
		Data: map[string]int64{
			"version": version,
		},
	})
}

func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - GetByEmail] "
	ctx := r.Context()

	vars := mux.Vars(r)
	email := vars["email"]

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

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - Update] "
	ctx := r.Context()

	if r.Method != http.MethodPut {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{})

	var requestData struct {
		User *models.User `json:"user"`
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	if requestData.User == nil {
		h.logger.Warn(ctx, ref+logger.LogMissingBodyData, map[string]any{})
		utils.ErrorResponse(w, fmt.Errorf("dados do usuário são obrigatórios"), http.StatusBadRequest)
		return
	}

	requestData.User.UID = id

	updatedUser, err := h.service.Update(ctx, requestData.User)
	if err != nil {
		if errors.Is(err, repository.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": updatedUser.UID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
		Data:    updatedUser,
	})
}

func (h *UserHandler) Disable(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - Disable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
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

	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version <= 0 {
		h.logger.Warn(ctx, ref+"versão inválida", map[string]any{
			"erro": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
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
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	user.Status = false
	user.Version = payload.Version

	_, err = h.service.Update(ctx, user)
	if err != nil {
		if errors.Is(err, repositories.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao desabilitar usuário: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Enable(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - Enable] "
	ctx := r.Context()

	if r.Method != http.MethodPatch {
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

	var payload struct {
		Version int `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version <= 0 {
		h.logger.Warn(ctx, ref+"versão inválida", map[string]any{
			"erro": err,
		})
		utils.ErrorResponse(w, fmt.Errorf("versão inválida"), http.StatusBadRequest)
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
			})
		}
		utils.ErrorResponse(w, err, status)
		return
	}

	user.Status = true
	user.Version = payload.Version

	_, err = h.service.Update(ctx, user)
	if err != nil {
		if errors.Is(err, repositories.ErrVersionConflict) {
			h.logger.Warn(ctx, ref+"conflito de versão", map[string]any{
				"user_id": id,
			})
			utils.ErrorResponse(w, fmt.Errorf("conflito de versão: os dados foram modificados por outro processo"), http.StatusConflict)
			return
		}
		h.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": id,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao habilitar usuário: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - Delete] "
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
