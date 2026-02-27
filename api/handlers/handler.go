// Package handlers contains HTTP request handlers for the KramerBot API.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/intothevoid/kramerbot/persist"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

// Handler holds the shared dependencies for all HTTP handlers.
type Handler struct {
	WebUserDB  persist.WebUserDBIF
	OzbScraper *scrapers.OzBargainScraper
	CCCScraper *scrapers.CamCamCamScraper
	Config     *util.Config
	Logger     *zap.Logger
	JWTSecret  []byte
}

// APIResponse is the standard JSON envelope returned by all endpoints.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func jsonCreated(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{Success: false, Error: msg})
}
