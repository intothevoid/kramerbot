package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/intothevoid/kramerbot/api/middleware"
	"go.uber.org/zap"
)

// GetProfile returns the authenticated user's profile (includes prefs and keywords).
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

// UpdatePreferences saves the user's deal notification toggles and syncs them
// to the bot's Telegram UserData record if the account is linked.
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

	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		h.Logger.Error("failed to fetch user for prefs update", zap.Error(err))
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	user.OzbGood = req.OzbGood
	user.OzbSuper = req.OzbSuper
	user.AmzDaily = req.AmzDaily
	user.AmzWeekly = req.AmzWeekly

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to save preferences", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.syncTelegramPrefs(user)

	jsonOK(w, user)
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

	kws := user.Keywords
	if kws == nil {
		kws = []string{}
	}
	jsonOK(w, map[string]interface{}{"keywords": kws})
}

type keywordRequest struct {
	Keyword string `json:"keyword"`
}

// AddKeyword appends a keyword to the user's watchlist and syncs to the bot.
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

	kw := strings.TrimSpace(strings.ToLower(req.Keyword))
	if kw == "" {
		jsonError(w, http.StatusBadRequest, "keyword cannot be empty")
		return
	}

	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	// Deduplicate.
	for _, existing := range user.Keywords {
		if existing == kw {
			jsonOK(w, map[string]interface{}{"keywords": user.Keywords})
			return
		}
	}

	user.Keywords = append(user.Keywords, kw)

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to save keyword", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.syncTelegramPrefs(user)

	jsonOK(w, map[string]interface{}{"keywords": user.Keywords})
}

// RemoveKeyword deletes a keyword from the user's watchlist and syncs to the bot.
func (h *Handler) RemoveKeyword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	kw := strings.ToLower(chi.URLParam(r, "keyword"))
	if kw == "" {
		jsonError(w, http.StatusBadRequest, "keyword param is required")
		return
	}

	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	filtered := user.Keywords[:0]
	for _, k := range user.Keywords {
		if k != kw {
			filtered = append(filtered, k)
		}
	}
	user.Keywords = filtered

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to remove keyword", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.syncTelegramPrefs(user)

	jsonOK(w, map[string]interface{}{"keywords": user.Keywords})
}
