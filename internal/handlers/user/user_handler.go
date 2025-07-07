package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
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
	if r.Method != http.MethodPost {
		h.logger.Warn(r.Context(), "[UserHandler] - Método não permitido", map[string]interface{}{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		User *models.User `json:"user"`
	}

	h.logger.Info(r.Context(), "[UserHandler] - Iniciando criação de usuário", map[string]interface{}{})

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(r.Context(), "[UserHandler] - Falha ao desserializar JSON", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.Create(r.Context(), requestData.User)
	if err != nil {
		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao criar usuário", map[string]interface{}{
			"email": requestData.User.Email,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Usuário criado com sucesso", map[string]interface{}{
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
	h.logger.Info(r.Context(), "[UserHandler] - Iniciando busca de todos usuários", nil)

	users, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao buscar usuários", nil)
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Usuários encontrados", map[string]interface{}{
		"quantidade": len(users),
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuários encontrados",
		Data:    users,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Context(), "[UserHandler] - Iniciando busca de usuário por ID", nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[UserHandler] - ID inválido na requisição", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
		}

		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao buscar usuário por ID", map[string]interface{}{
			"user_id": id,
			"status":  status,
		})

		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Usuário encontrado", map[string]interface{}{
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
	h.logger.Info(r.Context(), "[UserHandler] - Iniciando busca da versão do usuário por ID", nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[UserHandler] - ID inválido na requisição", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrUserNotFound) {
			status = http.StatusNotFound
		}
		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao obter versão do usuário", map[string]interface{}{
			"user_id": id,
			"status":  status,
		})
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Versão do usuário obtida com sucesso", map[string]interface{}{
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
	vars := mux.Vars(r)
	email := vars["email"]

	h.logger.Info(r.Context(), "[UserHandler] - Iniciando busca de usuário por email", map[string]interface{}{
		"email": email,
	})

	user, err := h.service.GetByEmail(r.Context(), email)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
		}
		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao buscar usuário por email", map[string]interface{}{
			"email":  email,
			"status": status,
		})
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Usuário encontrado por email", map[string]interface{}{
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
	if r.Method != http.MethodPut {
		h.logger.Warn(r.Context(), "[UserHandler] - Método não permitido", map[string]interface{}{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Iniciando atualização de usuário", nil)

	var requestData struct {
		User *models.User `json:"user"`
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[UserHandler] - ID inválido na requisição", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(r.Context(), "[UserHandler] - Falha ao desserializar JSON", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	if requestData.User == nil {
		h.logger.Warn(r.Context(), "[UserHandler] - Dados do usuário são obrigatórios", nil)
		utils.ErrorResponse(w, fmt.Errorf("dados do usuário são obrigatórios"), http.StatusBadRequest)
		return
	}

	requestData.User.UID = id

	updatedUser, err := h.service.Update(r.Context(), requestData.User)
	if err != nil {
		if errors.Is(err, repository.ErrVersionConflict) {
			h.logger.Warn(r.Context(), err.Error(), map[string]interface{}{
				"user_id": id,
			})
			utils.ErrorResponse(w, err, http.StatusConflict)
			return
		}
		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao atualizar usuário", map[string]interface{}{
			"user_id": id,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Usuário atualizado com sucesso", map[string]interface{}{
		"user_id": updatedUser.UID,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário atualizado com sucesso",
		Data:    updatedUser,
	})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.logger.Warn(r.Context(), "[UserHandler] - Método não permitido", map[string]interface{}{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Iniciando deleção de usuário", nil)

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		h.logger.Warn(r.Context(), "[UserHandler] - ID inválido na requisição", map[string]interface{}{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário não encontrado" {
			status = http.StatusNotFound
		}
		h.logger.Error(r.Context(), err, "[UserHandler] - Erro ao deletar usuário", map[string]interface{}{
			"user_id": id,
			"status":  status,
		})
		utils.ErrorResponse(w, err, status)
		return
	}

	h.logger.Info(r.Context(), "[UserHandler] - Usuário deletado com sucesso", map[string]interface{}{
		"user_id": id,
	})

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Usuário deletado com sucesso",
		Data:    nil,
	})
}
