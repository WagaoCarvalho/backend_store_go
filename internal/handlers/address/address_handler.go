package handlers

import (
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type AddressHandler struct {
	service services.AddressService
}

func NewAddressHandler(service services.AddressService) *AddressHandler {
	return &AddressHandler{service: service}
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var address models.Address
	if err := utils.FromJson(r.Body, &address); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdAddress, err := h.service.Create(r.Context(), address)
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

func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	address, err := h.service.GetByID(r.Context(), int(id))
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço encontrado",
		Data:    address,
	})
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	address.ID = &id

	if err := h.service.Update(r.Context(), address); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
	})
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), int(id)); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço deletado com sucesso",
	})
}
