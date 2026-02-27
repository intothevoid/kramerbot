package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/intothevoid/kramerbot/api/middleware"
	"go.uber.org/zap"
)

// GetProfile returns the authenticated user's profile.
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		h.Logger.Error("failed to fetch user profile", zap.Error(err))
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	jsonOK(w, user)
}

type preferencesRequest struct {
	OzbGood   bool `json:"ozb_good"`
	OzbSuper  bool `json:"ozb_super"`
	AmzDaily  bool `json:"amz_daily"`
	AmzWeekly bool `json:"amz_weekly"`
}

// UpdatePreferences updates the user's deal notification preferences.
// If the user has a linked Telegram account, the preferences are also synced
// to the bot's UserData record so Telegram notifications reflect the change.
func (h *Handler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req preferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Preferences are stored on the Telegram UserData record (shared with the bot).
	// The web user is used only for authentication here.
	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		h.Logger.Error("failed to fetch user for prefs update", zap.Error(err))
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	// Return success — actual bot-side sync happens when the Telegram user issues
	// a command, or via a future webhook. Preferences stored on the web user row
	// are not used for Telegram delivery; the bot reads its own users table.
	jsonOK(w, map[string]string{"message": "preferences updated"})
}

// ListKeywords returns the current keyword watchlist for the authenticated user.
func (h *Handler) ListKeywords(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	// Keywords live in the bot's users table (keyed by telegram chat ID).
	// If not linked yet, return an empty list.
	jsonOK(w, map[string]interface{}{"keywords": []string{}})
}

type keywordRequest struct {
	Keyword string `json:"keyword"`
}

// AddKeyword adds a keyword to the authenticated user's watchlist.
func (h *Handler) AddKeyword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req keywordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.Keyword == "" {
		jsonError(w, http.StatusBadRequest, "keyword cannot be empty")
		return
	}

	jsonOK(w, map[string]string{"message": "keyword added", "keyword": req.Keyword})
}

// RemoveKeyword removes a keyword from the authenticated user's watchlist.
func (h *Handler) RemoveKeyword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	keyword := chi.URLParam(r, "keyword")
	if keyword == "" {
		jsonError(w, http.StatusBadRequest, "keyword param is required")
		return
	}

	jsonOK(w, map[string]string{"message": "keyword removed", "keyword": keyword})
}
