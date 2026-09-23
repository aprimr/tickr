package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	er "github.com/aprimr/tickr/internal/pkg/errors"
	"github.com/aprimr/tickr/internal/pkg/response"
)

type AuthHandler interface {
	HandleUserRegister(w http.ResponseWriter, r *http.Request)
	HandleVenueAdminRegister(w http.ResponseWriter, r *http.Request)

	HandleVerifyAccount(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service AuthService
	log     *slog.Logger
}

func NewAuthHandler(service AuthService, logger *slog.Logger) AuthHandler {
	return &authHandler{
		service: service,
		log:     logger,
	}
}

// HandleUserRegister handles user registration
func (h *authHandler) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	var req UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode user registration request", "error", err)
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	userID, err := h.service.RegisterUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			response.Error(w, http.StatusConflict, ErrEmailAlreadyExists.Error(), map[string]string{"email": "an account with this email already exists"})
			return
		}

		h.log.Error("failed to register user", "error", err, "email", req.Email)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
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
		h.log.Warn("failed to decode user registration request", "error", err)
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	userID, err := h.service.RegisterVenueAdmin(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			response.Error(w, http.StatusConflict, ErrEmailAlreadyExists.Error(), map[string]string{"email": "an account with this email already exists"})
			return
		}

		h.log.Error("failed to register user", "error", err, "email", req.Email)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	res := map[string]any{
		"user_id": userID,
	}
	response.JSON(w, http.StatusCreated, "registration successful", res)
}

// HandleVerifyAccount handles the account verification request
func (h *authHandler) HandleVerifyAccount(w http.ResponseWriter, r *http.Request) {
	var req VerifyAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	err := h.service.VerifyUserAccount(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidOrExpiredOTP) {
			response.Error(w, http.StatusBadRequest, "invalid or expired OTP", nil)
			return
		}
		if errors.Is(err, ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found", nil)
			return
		}

		h.log.Error("failed to verify user account", "error", err, "user_id", req.UserID)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, "verification successful", nil)
}
