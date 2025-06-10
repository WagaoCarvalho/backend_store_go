package handlers

import (
	"errors"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
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

	createdAddress, err := h.service.Create(r.Context(), &address)
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

	address, err := h.service.GetByID(r.Context(), int64(id))
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

func (h *AddressHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	version, err := h.service.GetVersionByID(r.Context(), int64(id))
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Versão do endereço encontrada",
		Data:    map[string]int{"version": version},
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

	address.ID = id

	if err := h.service.Update(r.Context(), &address); err != nil {

		if ve, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(w, ve, http.StatusBadRequest)
			return
		}
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço atualizado com sucesso",
		Data:    nil,
	})
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.ErrorResponse(w, errors.New("ID inválido (esperado número inteiro)"), http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), int64(id))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrNotFound):
			utils.ErrorResponse(w, err, http.StatusNotFound)
		case errors.Is(err, services.ErrAddressIDRequired):
			utils.ErrorResponse(w, errors.New("endereço ID é obrigatório"), http.StatusBadRequest)
		default:
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Endereço deletado com sucesso",
		Data:    nil,
	})
}
