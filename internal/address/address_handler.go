package address

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aslam-ep/go-e-commerce/utils"
	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests related to user address.
type Handler struct {
	service Service
}

// NewHandler creates a new instance of the Handler with the provided address service.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) getIDsFromParam(r *http.Request) (int, int, error) {
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

// CreateAddress godoc
// @Summary      Adding a new address
// @Description  Adding a new address for the authenticated user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateUpdateAddressRequest  true  "Address request for create and update"
// @Success      200   {object}  Address
// @Failure      400   {object}  utils.MessageRes
// @Failure      401   {object}  utils.MessageRes
// @Failure      500   {object}  utils.MessageRes
// @Router       /users/{id}/addresses/create [post]
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var addressReq CreateUpdateAddressRequest
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

// GetAllAddress godoc
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
func (h *Handler) GetAllAddress(w http.ResponseWriter, r *http.Request) {
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

// GetAddressByID godoc
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
func (h *Handler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
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

// UpdateAddress godoc
// @Summary      Updated address by ID
// @Description  Updated address by ID for the current user
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path   int  true  "User ID"
// @Param        address_id     path   int  true  "Address ID"
// @Param        body  body CreateUpdateAddressRequest true  "Address request for create and update"
// @Success      200  {object}  Address
// @Failure      400  {object}  utils.MessageRes
// @Failure      401  {object}  utils.MessageRes
// @Failure      500  {object}  utils.MessageRes
// @Router       /users/{id}/addresses/{address_id}/update [put]
func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var addressReq CreateUpdateAddressRequest
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

// SetDefaultAddress godoc
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
func (h *Handler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
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

// DeleteAddress godoc
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
func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
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
