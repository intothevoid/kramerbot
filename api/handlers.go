package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/intothevoid/kramerbot/interfaces"
	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

// API holds dependencies for the API handlers
type API struct {
	BotBridge       interfaces.BotAPIBridge
	Logger          *zap.Logger
	DevelopmentMode bool // Whether development mode is enabled
}

// NewAPI creates a new API handler instance
func NewAPI(bridge interfaces.BotAPIBridge, logger *zap.Logger, developmentMode bool) *API {
	return &API{
		BotBridge:       bridge,
		Logger:          logger.Named("api"),
		DevelopmentMode: developmentMode,
	}
}

// Middleware to validate Telegram Web App initData
func (a *API) ValidateAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Logger.Info("Middleware called for: " + r.URL.Path)

		// Add CORS headers for local development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Telegram-Init-Data")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check if development mode is enabled
		if a.DevelopmentMode {
			// DEVELOPMENT MODE - Bypass authentication
			a.Logger.Info("Development mode enabled: Creating dummy user")
			dummyUser := &models.TelegramUser{
				ID:        12345, // Use any ID for testing
				FirstName: "Test",
				Username:  "testuser",
			}

			// Store user data in request context
			ctx := context.WithValue(r.Context(), userDataKey, dummyUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// PRODUCTION MODE - Normal authentication flow

		// Try to get initData from header first (normal Telegram WebApp flow)
		initDataString := r.Header.Get("X-Telegram-Init-Data")

		// If not in header, check if it's in URL query parameters (direct browser access)
		if initDataString == "" {
			// Check URL query for tgWebAppData parameter (our custom browser access)
			tgWebAppData := r.URL.Query().Get("tgWebAppData")
			if tgWebAppData != "" {
				a.Logger.Info("Using tgWebAppData from URL query parameter")

				// Extract chat_id from the tgWebAppData
				// Format is expected to be "chat_id=123456"
				if strings.HasPrefix(tgWebAppData, "chat_id=") {
					chatIDStr := strings.TrimPrefix(tgWebAppData, "chat_id=")
					chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
					if err == nil && chatID > 0 {
						a.Logger.Info("Creating user for direct browser access", zap.Int64("chatID", chatID))

						// Create a user based on the chat_id
						browserUser := &models.TelegramUser{
							ID:        chatID,
							FirstName: "Browser",
							Username:  "browser_user",
						}

						// Store user data in request context
						ctx := context.WithValue(r.Context(), userDataKey, browserUser)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			// No valid initData found in either header or URL
			a.Logger.Warn("Missing X-Telegram-Init-Data header and tgWebAppData parameter")
			http.Error(w, "Unauthorized: Missing authentication data", http.StatusUnauthorized)
			return
		}

		// Validate the initData if we have it
		userData, err := a.validateInitData(initDataString)
		if err != nil {
			a.Logger.Error("Invalid initData", zap.Error(err), zap.String("initData", initDataString))
			http.Error(w, fmt.Sprintf("Unauthorized: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		// Store user data in request context for downstream handlers
		ctx := context.WithValue(r.Context(), userDataKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateInitData checks the hash of the initData string
func (a *API) validateInitData(initData string) (*models.TelegramUser, error) {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("cannot parse initData query: %w", err)
	}

	hash := q.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash parameter not found in initData")
	}

	var checkData []string
	for k, v := range q {
		if k != "hash" {
			checkData = append(checkData, fmt.Sprintf("%s=%s", k, v[0]))
		}
	}
	sort.Strings(checkData)
	dataCheckString := strings.Join(checkData, "\n")

	secretKey := createSecretKey(a.BotBridge.GetBotToken())
	calculatedHash := calculateHMAC(secretKey, dataCheckString)

	if calculatedHash != hash {
		return nil, fmt.Errorf("invalid hash signature (calculated: %s, received: %s)", calculatedHash, hash)
	}

	// Optional: Check auth_date for expiration (e.g., within 24 hours)
	authDateUnix, err := strconv.ParseInt(q.Get("auth_date"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid auth_date: %w", err)
	}
	if time.Since(time.Unix(authDateUnix, 0)) > 24*time.Hour {
		return nil, fmt.Errorf("authentication data expired")
	}

	// Extract user data
	userJSON := q.Get("user")
	if userJSON == "" {
		return nil, fmt.Errorf("user data not found in initData")
	}

	var user models.TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, fmt.Errorf("cannot unmarshal user JSON: %w", err)
	}

	// Check if ID is present (crucial)
	if user.ID == 0 {
		return nil, fmt.Errorf("user ID is missing in initData")
	}

	return &user, nil
}

// Helper function to create the secret key for HMAC validation
func createSecretKey(botToken string) []byte {
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	return h.Sum(nil)
}

// Helper function to calculate HMAC-SHA256 hash
func calculateHMAC(secretKey []byte, data string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// --- Route Handlers ---

// HandlePreferences handles GET and POST requests for user preferences
func (a *API) HandlePreferences(w http.ResponseWriter, r *http.Request) {
	// Retrieve validated user from context
	userCtx := r.Context().Value(userDataKey)
	if userCtx == nil {
		a.Logger.Error("User data not found in context")
		http.Error(w, "Internal Server Error: User context missing", http.StatusInternalServerError)
		return
	}
	tgUser := userCtx.(*models.TelegramUser)
	chatID := tgUser.ID // UserID from Telegram is the ChatID for the bot

	a.Logger.Info("Handling preferences request", zap.String("method", r.Method), zap.Int64("chatID", chatID))

	switch r.Method {
	case http.MethodGet:
		a.getPreferences(w, r, chatID)
	case http.MethodPost:
		a.updatePreferences(w, r, chatID)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// getPreferences fetches and returns user preferences
func (a *API) getPreferences(w http.ResponseWriter, r *http.Request, chatID int64) {
	a.Logger.Info("getPreferences called", zap.Int64("chatID", chatID))

	user, err := a.BotBridge.GetUserDataWriter().GetUser(chatID)
	if err != nil {
		a.Logger.Error("Failed to get user preferences from DB",
			zap.Int64("chatID", chatID),
			zap.Error(err),
			zap.String("errorMessage", err.Error()))

		// Check if in development mode
		if a.DevelopmentMode {
			a.Logger.Info("Development mode: Returning dummy preferences for testing")

			// Create dummy user data for testing
			dummyUser := &models.UserData{
				ChatID:    chatID,
				Username:  "testuser",
				OzbGood:   true,
				OzbSuper:  false,
				AmzDaily:  true,
				AmzWeekly: false,
				Keywords:  []string{"laptop", "headphones", "monitor"},
				OzbSent:   []string{},
				AmzSent:   []string{},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(dummyUser); err != nil {
				a.Logger.Error("Failed to encode dummy preferences", zap.Error(err))
				http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
			} else {
				a.Logger.Info("Successfully sent dummy preferences")
			}
			return
		}

		// Production mode - return appropriate error
		if strings.Contains(err.Error(), "no rows") { // This check might need adjustment based on DB driver
			http.Error(w, "User not found. Please /start the bot first.", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve preferences", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		a.Logger.Error("Failed to encode preferences response", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
	} else {
		a.Logger.Info("Successfully sent real user preferences")
	}
}

// updatePreferences updates a single user preference
func (a *API) updatePreferences(w http.ResponseWriter, r *http.Request, chatID int64) {
	var reqBody map[string]bool
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.Logger.Error("Failed to decode update preference request body", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := a.BotBridge.GetUserDataWriter().GetUser(chatID)
	if err != nil {
		a.Logger.Error("Failed to get user for update", zap.Int64("chatID", chatID), zap.Error(err))

		// Check if in development mode
		if a.DevelopmentMode {
			a.Logger.Info("Development mode: Creating dummy user for testing")

			// Create dummy user data
			user = &models.UserData{
				ChatID:    chatID,
				Username:  "testuser",
				OzbGood:   reqBody["ozbGood"] || false,
				OzbSuper:  reqBody["ozbSuper"] || false,
				AmzDaily:  reqBody["amzDaily"] || false,
				AmzWeekly: reqBody["amzWeekly"] || false,
				Keywords:  []string{"laptop", "headphones", "monitor"},
				OzbSent:   []string{},
				AmzSent:   []string{},
			}

			// In development mode, just pretend the update succeeded
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"status": "success", "mode": "development"}`)
			return
		}

		// Production mode - return appropriate error
		if strings.Contains(err.Error(), "no rows") {
			http.Error(w, "User not found.", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
		}
		return
	}

	updated := false
	for key, value := range reqBody {
		switch key {
		case "ozbGood":
			if user.OzbGood != value {
				user.OzbGood = value
				updated = true
			}
		case "ozbSuper":
			if user.OzbSuper != value {
				user.OzbSuper = value
				updated = true
			}
		case "amzDaily":
			if user.AmzDaily != value {
				user.AmzDaily = value
				updated = true
			}
		case "amzWeekly":
			if user.AmzWeekly != value {
				user.AmzWeekly = value
				updated = true
			}
		default:
			a.Logger.Warn("Attempted to update unknown preference", zap.String("key", key), zap.Int64("chatID", chatID))
			// Maybe return an error here if only specific keys are allowed?
			// For now, just ignore unknown keys.
		}
	}

	if !updated {
		w.WriteHeader(http.StatusOK) // Nothing changed, but request was valid
		fmt.Fprintln(w, `{"status": "no changes"}`)
		return
	}

	// Update in DB and in-memory store using the bridge method
	if err := a.BotBridge.UpdateUserData(user); err != nil {
		a.Logger.Error("Failed to update user preferences via bridge", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Failed to save preferences", http.StatusInternalServerError)
		return
	}

	// No longer need to update memory store here, bridge method handles it
	// a.Bot.UserStore.Users[chatID] = user

	a.Logger.Info("Successfully updated preferences", zap.Int64("chatID", chatID))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "success"}`) // Send a simple success status
}

// --- Keyword Handlers ---

type KeywordRequest struct {
	Keyword string `json:"keyword"`
}

type KeywordsResponse struct {
	Keywords []string `json:"keywords"`
}

// HandleAddKeyword adds a keyword for the user
func (a *API) HandleAddKeyword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userCtx := r.Context().Value(userDataKey)
	tgUser := userCtx.(*models.TelegramUser)
	chatID := tgUser.ID

	var req KeywordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.Logger.Error("Failed to decode add keyword request body", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keyword := strings.TrimSpace(strings.ToLower(req.Keyword))
	if keyword == "" {
		http.Error(w, "Keyword cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := a.BotBridge.GetUserDataWriter().GetUser(chatID)
	if err != nil {
		a.Logger.Error("Failed to get user for adding keyword", zap.Int64("chatID", chatID), zap.Error(err))

		// Check if in development mode
		if a.DevelopmentMode {
			a.Logger.Info("Development mode: Using dummy user for keyword addition")

			// Create dummy response with the new keyword
			keywords := []string{"laptop", "headphones", "monitor", keyword}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(KeywordsResponse{Keywords: keywords})
			return
		}

		http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
		return
	}

	// Check if keyword already exists
	found := false
	for _, existing := range user.Keywords {
		if existing == keyword {
			found = true
			break
		}
	}
	if found {
		a.Logger.Info("Keyword already exists", zap.String("keyword", keyword), zap.Int64("chatID", chatID))
		// Return current list even if keyword exists (idempotent-like)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(KeywordsResponse{Keywords: user.Keywords})
		return
	}

	// Add keyword
	user.Keywords = append(user.Keywords, keyword)

	// Update DB and memory via bridge
	if err := a.BotBridge.UpdateUserData(user); err != nil {
		a.Logger.Error("Failed to update user keywords via bridge (add)", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Failed to add keyword", http.StatusInternalServerError)
		return
	}

	a.Logger.Info("Successfully added keyword", zap.String("keyword", keyword), zap.Int64("chatID", chatID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(KeywordsResponse{Keywords: user.Keywords}) // Return updated list
}

// HandleRemoveKeyword removes a keyword for the user
func (a *API) HandleRemoveKeyword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userCtx := r.Context().Value(userDataKey)
	tgUser := userCtx.(*models.TelegramUser)
	chatID := tgUser.ID

	var req KeywordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.Logger.Error("Failed to decode remove keyword request body", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keywordToRemove := strings.TrimSpace(strings.ToLower(req.Keyword))
	if keywordToRemove == "" {
		http.Error(w, "Keyword cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := a.BotBridge.GetUserDataWriter().GetUser(chatID)
	if err != nil {
		a.Logger.Error("Failed to get user for removing keyword", zap.Int64("chatID", chatID), zap.Error(err))

		// Check if in development mode
		if a.DevelopmentMode {
			a.Logger.Info("Development mode: Using dummy user for keyword removal")

			// Create dummy keywords list without the keyword to remove
			dummyKeywords := []string{"laptop", "headphones", "monitor"}
			var updatedKeywords []string

			for _, k := range dummyKeywords {
				if k != keywordToRemove {
					updatedKeywords = append(updatedKeywords, k)
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(KeywordsResponse{Keywords: updatedKeywords})
			return
		}

		http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
		return
	}

	found := false
	var updatedKeywords []string
	for _, existing := range user.Keywords {
		if existing != keywordToRemove {
			updatedKeywords = append(updatedKeywords, existing)
		} else {
			found = true
		}
	}

	if !found {
		a.Logger.Info("Keyword not found for removal", zap.String("keyword", keywordToRemove), zap.Int64("chatID", chatID))
		// Return current list even if keyword wasn't found (idempotent-like)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(KeywordsResponse{Keywords: user.Keywords})
		return
	}

	// Update keywords
	user.Keywords = updatedKeywords

	// Update DB and memory via bridge
	if err := a.BotBridge.UpdateUserData(user); err != nil {
		a.Logger.Error("Failed to update user keywords via bridge (remove)", zap.Int64("chatID", chatID), zap.Error(err))
		http.Error(w, "Failed to remove keyword", http.StatusInternalServerError)
		return
	}
	// No longer need to update memory store here
	// a.Bot.UserStore.Users[chatID].Keywords = user.Keywords

	a.Logger.Info("Successfully removed keyword", zap.String("keyword", keywordToRemove), zap.Int64("chatID", chatID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(KeywordsResponse{Keywords: user.Keywords}) // Return updated list
}

// --- Action Handlers ---

// HandleTestNotification sends a test notification to the user
func (a *API) HandleTestNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userCtx := r.Context().Value(userDataKey)
	tgUser := userCtx.(*models.TelegramUser)
	chatID := tgUser.ID

	a.Logger.Info("Received request to send test notification", zap.Int64("chatID", chatID))

	// Check if in development mode
	if a.DevelopmentMode || chatID == 12345 {
		a.Logger.Info("Development mode: Pretending to send test notification")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "success", "message": "Test notification sent (development mode)"}`)
		return
	}

	// Call the bridge method for real users
	if err := a.BotBridge.SendTestMessageToChat(chatID, tgUser); err != nil {
		a.Logger.Error("Failed to send test notification via bridge", zap.Int64("chatID", chatID), zap.Error(err))
		// Decide if this error should be shown to the user
		http.Error(w, "Failed to trigger test notification", http.StatusInternalServerError)
		return
	}

	a.Logger.Info("Successfully triggered test notification", zap.Int64("chatID", chatID))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "success", "message": "Test notification sent"}`)
}

// --- Context Key ---
// It's good practice to define a custom type for context keys
type contextKey string

const userDataKey contextKey = "userData"
