package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/gorilla/mux"
)

type AddressHandler struct {
	service services.AddressService
}

func NewAddressHandler(service services.AddressService) *AddressHandler {
	return &AddressHandler{service: service}
}

// Criar um endereço
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		http.Error(w, `{"status":400, "message":"Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	createdAddress, err := h.service.CreateAddress(r.Context(), address)
	if err != nil {
		http.Error(w, `{"status":500, "message":"Erro ao criar endereço"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "Endereço criado com sucesso",
		"data":    createdAddress,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Buscar endereço por ID
func (h *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, `{"status":400, "message":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	address, err := h.service.GetAddressByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"status":404, "message":"Endereço não encontrado"}`, http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Endereço encontrado",
		"data":    address,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Atualizar endereço
func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, `{"status":400, "message":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		http.Error(w, `{"status":400, "message":"Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	// Atribuir o ID corretamente
	address.ID = id

	if err := h.service.UpdateAddress(r.Context(), address); err != nil {
		http.Error(w, `{"status":500, "message":"Erro ao atualizar endereço"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Endereço atualizado com sucesso",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Deletar endereço
func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, `{"status":400, "message":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAddress(r.Context(), id); err != nil {
		http.Error(w, `{"status":500, "message":"Erro ao deletar endereço"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Endereço deletado com sucesso",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
