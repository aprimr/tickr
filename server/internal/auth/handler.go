package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/aprimr/tickr/internal/domain"
	"github.com/aprimr/tickr/internal/middleware"
	"github.com/aprimr/tickr/internal/pkg/device"
	er "github.com/aprimr/tickr/internal/pkg/errors"
	"github.com/aprimr/tickr/internal/pkg/response"
	"github.com/aprimr/tickr/internal/utils/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)

	HandleUserRegister(w http.ResponseWriter, r *http.Request)
	HandleVenueAdminRegister(w http.ResponseWriter, r *http.Request)

	HandleVerifyAccount(w http.ResponseWriter, r *http.Request)
	HandleForgotPassword(w http.ResponseWriter, r *http.Request)
	HandleResetPassword(w http.ResponseWriter, r *http.Request)

	HandleTokenRotation(w http.ResponseWriter, r *http.Request)

	HandleLogout(w http.ResponseWriter, r *http.Request)
	HandleLogoutAllDevices(w http.ResponseWriter, r *http.Request)
	HandleGetActiveSessions(w http.ResponseWriter, r *http.Request)
	HandleDeleteSession(w http.ResponseWriter, r *http.Request)
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

// HandleLogin handles user authentication for user, venue admin and super admin
func (h *authHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode login request", "error", err)
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	// Get device info from request
	deviceInfo := device.Parse(r)

	accessToken, refreshToken, err := h.service.Login(r.Context(), req, deviceInfo)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, ErrInvalidCredentials.Error(), nil)
			return
		}
		if errors.Is(err, ErrEmailNotVerified) {
			response.Error(w, http.StatusForbidden, ErrEmailNotVerified.Error(), nil)
			return
		}
		if errors.Is(err, ErrAccountDeactivated) {
			response.Error(w, http.StatusForbidden, ErrAccountDeactivated.Error(), nil)
			return
		}

		h.log.Error("login failed", "error", err, "email", req.Email)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	res := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.JSON(w, http.StatusOK, "login successful", res)
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

	res := RegisterResponse{
		UserID: userID,
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

	res := RegisterResponse{
		UserID: userID,
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

// HandleForgotPassword handles the request for forget password
func (h *authHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	err := h.service.ForgotPassword(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrActiveOTPAlreadyExists) {
			response.Error(w, http.StatusTooManyRequests, ErrActiveOTPAlreadyExists.Error(), nil)
			return
		}

		h.log.Error("failed to process forget password request", "error", err, "email", req.Email)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, "password reset code sent to the email", nil)
}

// HandleResetPassword handles the password resetting request
func (h *authHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	err := h.service.ResetPassword(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidOrExpiredOTP) {
			response.Error(w, http.StatusBadRequest, ErrInvalidOrExpiredOTP.Error(), nil)
			return
		}

		h.log.Error("failed to reset user password", "error", err, "email", req.Email)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, "password reset successful", nil)
}

// HandleTokenRotation handles the rotation of the tokens
func (h *authHandler) HandleTokenRotation(w http.ResponseWriter, r *http.Request) {
	var req RotateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	// Get device info
	deviceInfo := device.Parse(r)

	accessToken, refreshToken, err := h.service.RotateToken(r.Context(), req, deviceInfo)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrExpiredToken):
			response.Error(w, http.StatusUnauthorized, jwt.ErrExpiredToken.Error(), nil)

		case errors.Is(err, jwt.ErrInvalidToken):
			response.Error(w, http.StatusUnauthorized, jwt.ErrInvalidToken.Error(), nil)

		case errors.Is(err, domain.ErrRefreshTokenNotFound):
			response.Error(w, http.StatusUnauthorized, domain.ErrRefreshTokenNotFound.Error(), nil)

		default:
			response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		}

		h.log.Error("failed to rotate refresh token", "error", err)
		return
	}

	res := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.JSON(w, http.StatusOK, "token rotation successful", res)
}

// HandleLogout logs out user by revoking the refresh token session
func (h *authHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, er.ErrInvalidReqBody.Error(), nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, "validation failed", validationErr)
		return
	}

	err := h.service.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Error("logout failed", "error", err)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, "logout successful", nil)
}

// HandleLogoutAllDevices handles revoking all sessions for the current user
func (h *authHandler) HandleLogoutAllDevices(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err := h.service.LogoutAllDevices(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		h.log.Error("logout all device failed", "error", err, "user_id", userID)
		return
	}

	response.JSON(w, http.StatusOK, "logout successful", nil)
}

// HandleGetActiveSessions returns all active sessions for the logged-in user
func (h *authHandler) HandleGetActiveSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Get access token from header `x-refresh-token`
	refreshToken := r.Header.Get("x-refresh-token")
	if refreshToken == "" {
		response.Error(w, http.StatusUnauthorized, "missing refresh token header", nil)
		return
	}

	sessions, err := h.service.GetActiveSessions(r.Context(), userID, refreshToken)
	if err != nil {
		h.log.Error("get active sessions failed", "error", err, "user_id", userID)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, "sessions fetched successfully", sessions)
}

// HandleDeleteSession deletes a specific session by its ID
func (h *authHandler) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Get session id from URL params
	// (e.g., /auth/sessions/123-abc)
	sessionIDStr := chi.URLParam(r, "sessionID")
	if sessionIDStr == "" {
		response.Error(w, http.StatusBadRequest, "session id is required", nil)
		return
	}

	// Convert SessionIdStr to UUID
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid session id format", nil)
		return
	}

	err = h.service.DeleteSession(r.Context(), userID, sessionID)
	if err != nil {
		h.log.Error("delete session failed", "error", err, "user_id", userID, "session_id", sessionID)
		response.Error(w, http.StatusInternalServerError, er.ErrInternalError.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, "session delete successful", nil)
}
