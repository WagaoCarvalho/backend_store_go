package handlers

import (
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type AddressHandlerInterface interface {
	CreateAddress(w http.ResponseWriter, r *http.Request)
	GetAddress(w http.ResponseWriter, r *http.Request)
	UpdateAddress(w http.ResponseWriter, r *http.Request)
	DeleteAddress(w http.ResponseWriter, r *http.Request)
}

type AddressHandler struct {
	service services.AddressService
}

func NewAddressHandler(service services.AddressService) AddressHandlerInterface {
	return &AddressHandler{service: service}
}

// Criar um endereço
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address models.Address
	if err := utils.FromJson(r.Body, &address); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdAddress, err := h.service.CreateAddress(r.Context(), address)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Endereço criado com sucesso",
		Data:    createdAddress,
	})
}

// Buscar endereço por ID
func (h *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	address, err := h.service.GetAddressByID(r.Context(), int(id))
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Endereço criado com sucesso",
		Data:    address,
	})
}

// Atualizar endereço
func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var address models.Address
	if err := utils.FromJson(r.Body, &address); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	address.ID = int(id)

	if err := h.service.UpdateAddress(r.Context(), address); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
	})
}

// Deletar endereço
func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAddress(r.Context(), int(id)); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço deletado com sucesso",
	})
}
