package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/intothevoid/kramerbot/api/middleware"
	"go.uber.org/zap"
)

type telegramLinkResponse struct {
	Token    string `json:"token"`
	DeepLink string `json:"deep_link"`
	ExpiresAt time.Time `json:"expires_at"`
}

type telegramStatusResponse struct {
	Linked           bool    `json:"linked"`
	TelegramUsername *string `json:"telegram_username,omitempty"`
}

// GenerateTelegramLink creates a one-time deep link token so the user can connect
// their Telegram account by clicking t.me/<botUsername>?start=<token>.
func (h *Handler) GenerateTelegramLink(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.WebUserDB.GetWebUserByID(claims.UserID)
	if err != nil || user == nil {
		h.Logger.Error("failed to fetch user for telegram link", zap.Error(err))
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	token := uuid.New().String()
	expires := time.Now().UTC().Add(15 * time.Minute)
	user.LinkToken = &token
	user.LinkTokenExpires = &expires

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to save link token", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	botUsername := os.Getenv("TELEGRAM_BOT_USERNAME")
	if botUsername == "" {
		botUsername = "kramerbot"
	}
	deepLink := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, token)

	jsonOK(w, telegramLinkResponse{
		Token:    token,
		DeepLink: deepLink,
		ExpiresAt: expires,
	})
}

// GetTelegramStatus returns whether the authenticated user has a linked Telegram account.
func (h *Handler) GetTelegramStatus(w http.ResponseWriter, r *http.Request) {
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

	resp := telegramStatusResponse{
		Linked:           user.TelegramChatID != nil,
		TelegramUsername: user.TelegramUsername,
	}
	jsonOK(w, resp)
}

// UnlinkTelegram removes the Telegram account association from the web user.
func (h *Handler) UnlinkTelegram(w http.ResponseWriter, r *http.Request) {
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

	user.TelegramChatID = nil
	user.TelegramUsername = nil
	user.LinkToken = nil
	user.LinkTokenExpires = nil

	if err := h.WebUserDB.UpdateWebUser(user); err != nil {
		h.Logger.Error("failed to unlink telegram", zap.Error(err))
		jsonError(w, http.StatusInternalServerError, "internal error")
		return
	}

	jsonOK(w, map[string]string{"message": "telegram unlinked"})
}
