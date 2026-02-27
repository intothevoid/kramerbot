// Package api provides the HTTP API server for the KramerBot web interface.
package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/intothevoid/kramerbot/api/handlers"
	"github.com/intothevoid/kramerbot/api/middleware"
	"github.com/intothevoid/kramerbot/persist"
	sqlite_persist "github.com/intothevoid/kramerbot/persist/sqlite"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

// Server is the HTTP API server.
type Server struct {
	httpServer *http.Server
	Logger     *zap.Logger
}

// NewServer constructs the Chi router, registers all routes, and returns a ready-to-run Server.
func NewServer(
	cfg *util.Config,
	db persist.DatabaseIF,
	ozbScraper *scrapers.OzBargainScraper,
	cccScraper *scrapers.CamCamCamScraper,
	logger *zap.Logger,
) (*Server, error) {
	// Resolve JWT secret — prefer env var over a generated fallback.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Warn("JWT_SECRET env var not set; API authentication will be insecure in production")
		jwtSecret = "changeme-set-JWT_SECRET-in-production"
	}

	// Cast to WebUserDBIF so handlers can manage web_users.
	webUserDB, ok := db.(persist.WebUserDBIF)
	if !ok {
		// Try via the SQLiteWrapper directly.
		if sw, ok2 := db.(*sqlite_persist.SQLiteWrapper); ok2 {
			webUserDB = sw
		} else {
			return nil, fmt.Errorf("database driver does not implement WebUserDBIF")
		}
	}

	h := &handlers.Handler{
		WebUserDB:  webUserDB,
		OzbScraper: ozbScraper,
		CCCScraper: cccScraper,
		Config:     cfg,
		Logger:     logger,
		JWTSecret:  []byte(jwtSecret),
	}

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(middleware.CORS(cfg.API.CORSOrigins))

	// Public routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)
	})

	// Protected routes
	r.Route("/api/v1/user", func(r chi.Router) {
		r.Use(middleware.JWTAuth([]byte(jwtSecret)))
		r.Get("/profile", h.GetProfile)
		r.Put("/preferences", h.UpdatePreferences)
		r.Get("/keywords", h.ListKeywords)
		r.Post("/keywords", h.AddKeyword)
		r.Delete("/keywords/{keyword}", h.RemoveKeyword)
		r.Post("/telegram/link", h.GenerateTelegramLink)
		r.Get("/telegram/status", h.GetTelegramStatus)
		r.Delete("/telegram/link", h.UnlinkTelegram)
	})

	// Deal feed (requires auth)
	r.Route("/api/v1/deals", func(r chi.Router) {
		r.Use(middleware.JWTAuth([]byte(jwtSecret)))
		r.Get("/ozbargain", h.GetOzbDeals)
		r.Get("/amazon", h.GetAmazonDeals)
		r.Get("/", h.GetAllDeals)
	})

	// Health check (public)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	srv := &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Logger: logger,
	}

	return srv, nil
}

// Start begins listening on the configured address. It blocks until the server exits.
func (s *Server) Start() error {
	s.Logger.Info("API server listening", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
