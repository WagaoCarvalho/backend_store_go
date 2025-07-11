package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/login"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

type LoginHandler struct {
	service auth.LoginService
}

func NewLoginHandler(service auth.LoginService) *LoginHandler {
	return &LoginHandler{service: service}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var credentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	response := utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Login realizado com sucesso",
		Data:    map[string]string{"token": token},
	}
	utils.ToJson(w, http.StatusOK, response)
}
