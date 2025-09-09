package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service service.UserService
	logger  *logger.LogAdapter
}

func NewUserHandler(service service.UserService, logger *logger.LogAdapter) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	// Recebe o DTO no lugar do model
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

	// Converte DTO para model antes de enviar para o service
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

	// Converte model criado de volta para DTO antes de retornar
	createdDTO := dto.ToUserDTO(createdUser)

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdDTO,
	})
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    users,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	})
}

func (h *UserHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetVersionByID] "
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
		if errors.Is(err, errMsg.ErrNotFound) {
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

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do usuário obtida com sucesso",
		Data: map[string]int64{
			"version": version,
		},
	})
}

func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetByEmail] "
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

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário encontrado",
		Data:    user,
	})
}

func (h *UserHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - GetByName] "
	ctx := r.Context()

	vars := mux.Vars(r)
	name := vars["username"]

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

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    users,
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Update] "
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
		User *dto.UserDTO `json:"user"` // <-- agora DTO
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
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

	// Converte DTO para model antes de enviar para o service
	userModel := dto.ToUserModel(*requestData.User)
	userModel.UID = id

	updatedUser, err := h.service.Update(ctx, userModel)
	if err != nil {
		if errors.Is(err, errMsg.ErrVersionConflict) {
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

	// Converte model de volta para DTO na resposta
	updatedDTO := dto.ToUserDTO(updatedUser)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
		Data:    updatedDTO,
	})
}

func (h *UserHandler) Disable(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - Disable] "
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
		if errors.Is(err, errMsg.ErrVersionConflict) {
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
	const ref = "[UserHandler - Enable] "
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
		if errors.Is(err, errMsg.ErrVersionConflict) {
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
