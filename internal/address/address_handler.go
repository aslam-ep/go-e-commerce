package address

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aslam-ep/go-e-commerce/utils"
	"github.com/go-chi/chi/v5"
)

type AddressHandler struct {
	service AddressService
}

func NewAddressHandler(addressService AddressService) *AddressHandler {
	return &AddressHandler{
		service: addressService,
	}
}

func (h *AddressHandler) getIDsFromParam(r *http.Request) (int, int, error) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return -1, -1, err
	}

	addressIDStr := chi.URLParam(r, "address_id")
	addressID, err := strconv.Atoi(addressIDStr)
	if err != nil {
		return -1, -1, err
	}

	return addressID, userID, nil
}

// @Summary      Adding a new address
// @Description  Adding a new address for the authenticated user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      AddressReq  true  "Address request for create and update"
// @Success      200   {object}  Address
// @Failure      400   {object}  utils.MessageRes
// @Failure      401   {object}  utils.MessageRes
// @Failure      500   {object}  utils.MessageRes
// @Router       /users/{id}/addresses/create [post]
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var addressReq AddressReq
	if err := utils.ReadFromRequest(r, &addressReq); err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("id: ", err)
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	addressReq.UserID = int64(userID)

	if err := utils.Validate.Struct(addressReq); err != nil {
		log.Println("validation: ", err)
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateAddress(r.Context(), &addressReq)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteResponse(w, http.StatusOK, res)
}

// @Summary      Get all addresses
// @Description  Get all addresses for the authenticated user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  AddressRes
// @Param        id             path   int  true  "User ID"
// @Failure      400  {object}  utils.MessageRes
// @Failure      401   {object}  utils.MessageRes
// @Failure      404  {object}  utils.MessageRes
// @Router       /users/{id}/addresses/ [get]
func (h *AddressHandler) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.GetAllAddress(r.Context(), userID)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteResponse(w, http.StatusOK, res)
}

// @Summary      Get address by ID
// @Description  Get a specific address by ID for the authenticated user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path   int  true  "User ID"
// @Param        address_id     path   int  true  "Address ID"
// @Success      200  {object}  Address
// @Failure      400  {object}  utils.MessageRes
// @Failure      401  {object}  utils.MessageRes
// @Failure      404  {object}  utils.MessageRes
// @Router       /users/{id}/addresses/{adddress_id} [get]
func (h *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	addressID, userID, err := h.getIDsFromParam(r)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.GetAddressByID(r.Context(), addressID, userID)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteResponse(w, http.StatusOK, res)
}

// @Summary      Updated address by ID
// @Description  Updated address by ID for the current user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path   int  true  "User ID"
// @Param        address_id     path   int  true  "Address ID"
// @Param        body  body AddressReq true  "Address request for create and update"
// @Success      200  {object}  Address
// @Failure      400  {object}  utils.MessageRes
// @Failure      401  {object}  utils.MessageRes
// @Failure      500  {object}  utils.MessageRes
// @Router       /users/{id}/addresses/{address_id}/update [put]
func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var addressReq AddressReq
	if err := utils.ReadFromRequest(r, &addressReq); err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	addressID, userID, err := h.getIDsFromParam(r)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	addressReq.ID = int64(addressID)
	addressReq.UserID = int64(userID)

	if err := utils.Validate.Struct(addressReq); err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateAddress(r.Context(), &addressReq)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteResponse(w, http.StatusOK, res)
}

// @Summary      Set default address
// @Description  Set a specific address as the default for the authenticated user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path   int  true  "User ID"
// @Param        address_id     path   int  true  "Address ID"
// @Success      200  {object}  utils.MessageRes
// @Failure      400  {object}  utils.MessageRes
// @Failure      401  {object}  utils.MessageRes
// @Failure      404  {object}  utils.MessageRes
// @Router       /users/{id}/addresses/{address_id}/set-default [put]
func (h *AddressHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	addressID, userID, err := h.getIDsFromParam(r)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.SetDefaultAddress(r.Context(), addressID, userID)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteResponse(w, http.StatusOK, res)
}

// @Summary      Delete address
// @Description  Delete a specific address by ID for the authenticated user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path   int  true  "User ID"
// @Param        address_id     path   int  true  "Address ID"
// @Success      200  {object}  utils.MessageRes
// @Failure      400  {object}  utils.MessageRes
// @Failure      401  {object}  utils.MessageRes
// @Failure      404  {object}  utils.MessageRes
// @Router       /users/{id}/addresses/{address_id}/delete [delete]
func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressID, userID, err := h.getIDsFromParam(r)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.DeleteAddress(r.Context(), addressID, userID)
	if err != nil {
		utils.WriterErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteResponse(w, http.StatusOK, res)
}
