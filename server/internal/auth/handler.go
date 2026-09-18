package auth

import (
	"encoding/json"
	"net/http"

	"github.com/aprimr/tickr/internal/pkg/response"
)

type AuthHandler interface {
	HandleUserRegister(w http.ResponseWriter, r *http.Request)
	HandleVenueAdminRegister(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) AuthHandler {
	return &authHandler{service: service}
}

// HandleUserRegister handles user registration
func (h *authHandler) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	var req UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if validationErr := req.Validate(); validationErr != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	userID, err := h.service.RegisterUser(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	res := map[string]any{
		"user_id": userID,
	}
	response.JSON(w, http.StatusCreated, "registration successful", res)
}

// HandleVenueAdminRegister handles venue admin registration
func (h *authHandler) HandleVenueAdminRegister(w http.ResponseWriter, r *http.Request) {
	var req VenueRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if validationErr := req.Validate(); validationErr != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	userID, err := h.service.RegisterVenueAdmin(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	res := map[string]any{
		"user_id": userID,
	}
	response.JSON(w, http.StatusCreated, "registration successful", res)
}
