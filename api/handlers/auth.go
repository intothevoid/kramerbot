package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/intothevoid/kramerbot/api/middleware"
	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string          `json:"token"`
	User  *models.WebUser `json:"user"`
}

// Register creates a new web user account.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.Email == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "email and password are required")
		return
	}
	if len(req.Password) < 8 {
		jsonError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	existing, err := h.WebUserDB.GetWebUserByEmail(req.Email)
	if err != nil {
		h.Logger.Error("failed to check existing user", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if existing != nil {
		jsonError(w, http.StatusConflict, "an account with this email already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("failed to hash password", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	user := &models.WebUser{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
		DisplayName:  req.DisplayName,
		Keywords:     []string{},
	}
	if err := h.WebUserDB.CreateWebUser(user); err != nil {
		h.Logger.Error("failed to create web user", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	token, err := h.signToken(user)
	if err != nil {
		h.Logger.Error("failed to sign JWT", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	jsonCreated(w, authResponse{Token: token, User: user})
}

// Login authenticates a web user and returns a JWT.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := h.WebUserDB.GetWebUserByEmail(req.Email)
	if err != nil {
		h.Logger.Error("failed to look up user", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		jsonError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := h.signToken(user)
	if err != nil {
		h.Logger.Error("failed to sign JWT", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	jsonOK(w, authResponse{Token: token, User: user})
}

// Logout is a no-op for stateless JWTs — the client discards the token.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]string{"message": "logged out"})
}

// ForgotPassword generates a password reset token and returns the reset link.
// In production, you would email this link. Here we return it in the response
// and log it — suitable for a self-hosted personal bot.
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" {
		jsonError(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := h.WebUserDB.GetWebUserByEmail(req.Email)
	if err != nil {
		h.Logger.Error("forgot-password lookup error", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Always return the same generic message to prevent email enumeration.
	generic := map[string]string{"message": "If an account exists for that email, a reset link has been generated."}

	if user == nil {
		jsonOK(w, generic)
		return
	}

	token := uuid.New().String()
	expires := time.Now().UTC().Add(1 * time.Hour)
	user.ResetToken = &token
	user.ResetTokenExpires = &expires

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to save reset token", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Build the reset link. The frontend base URL is taken from CORS origins or an env var.
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" && len(h.Config.API.CORSOrigins) > 0 {
		frontendURL = h.Config.API.CORSOrigins[0]
	}
	if frontendURL == "" {
		frontendURL = fmt.Sprintf("http://localhost:%d", h.Config.API.Port)
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	h.Logger.Info("password reset token generated",
		zap.String("email", user.Email),
		zap.String("reset_link", resetLink),
	)

	// Return the link directly (for self-hosted/personal use without SMTP).
	jsonOK(w, map[string]string{
		"message":    "Reset link generated. Copy it below (in production this would be emailed).",
		"reset_link": resetLink,
	})
}

// ResetPassword validates a reset token and updates the password.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Token == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "token and password are required")
		return
	}
	if len(req.Password) < 8 {
		jsonError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	user, err := h.WebUserDB.GetWebUserByResetToken(req.Token)
	if err != nil {
		h.Logger.Error("reset-password lookup error", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if user == nil {
		jsonError(w, http.StatusBadRequest, "invalid or expired reset token")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("failed to hash password", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	user.PasswordHash = string(hash)
	user.ResetToken = nil
	user.ResetTokenExpires = nil

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to update password", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	jsonOK(w, map[string]string{"message": "password updated successfully"})
}

// signToken creates a signed JWT for the given web user.
func (h *Handler) signToken(user *models.WebUser) (string, error) {
	expiry := time.Duration(h.Config.API.JWTExpiryHours) * time.Hour
	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.JWTSecret)
}
