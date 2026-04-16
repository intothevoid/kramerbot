package main

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/intothevoid/kramerbot/api"
	"github.com/intothevoid/kramerbot/bot"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	sqlite_persist "github.com/intothevoid/kramerbot/persist/sqlite"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// create a new instance of our bot
	k := new(bot.KramerBot)

	// Setup configuration
	confPath, _ := os.Getwd()
	confPath += "/config.yaml"

	// Initialize logger first with default settings
	logger := util.SetupLogger(zapcore.DebugLevel, true)

	// Load configuration
	config, err := util.SetupConfig(confPath, logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Update logger with config settings
	logger = util.SetupLogger(zapcore.Level(config.LogLevel), config.LogToFile)
	k.Logger = logger
	k.Config = config

	// Android TV notifications via Pipup
	if config.Pipup.Enabled {
		k.Pipup = pipup.New(config.Pipup, logger)
	}

	// Get the token for the telegram bot api
	k.Token = k.GetToken()

	// Test mode doesn't require a token
	if k.Token == "" && !config.TestMode {
		logger.Fatal("Cannot proceed without a bot token, is the TELEGRAM_BOT_TOKEN environment variable set?")
	}

	// Create Ozbargain scraper
	ozbscraper := new(scrapers.OzBargainScraper)
	ozbscraper.SID = scrapers.SID_OZBARGAIN
	ozbscraper.Logger = logger
	ozbscraper.BaseUrl = scrapers.URL_OZBARGAIN
	ozbscraper.Deals = []models.OzBargainDeal{}
	ozbscraper.ScrapeInterval = config.Scrapers.OzBargain.ScrapeInterval
	ozbscraper.MaxDealsToStore = config.Scrapers.OzBargain.MaxStoredDeals

	// Create CamelCamelCamel (Amazon) scraper
	cccscraper := new(scrapers.CamCamCamScraper)
	cccscraper.SID = scrapers.SID_CCC_AMAZON
	cccscraper.Logger = logger
	cccscraper.BaseUrl = config.Scrapers.Amazon.URLs
	cccscraper.Deals = []models.CamCamCamDeal{}
	cccscraper.ScrapeInterval = config.Scrapers.Amazon.ScrapeInterval
	cccscraper.MaxDealsToStore = config.Scrapers.Amazon.MaxStoredDeals

	// Initialise bot (creates DB connection internally)
	k.NewBot(ozbscraper, cccscraper)

	// Wire the WebUserDB so the bot can resolve Telegram link tokens.
	if sw, ok := k.DataWriter.(*sqlite_persist.SQLiteWrapper); ok {
		k.WebUserDB = sw
	} else {
		logger.Warn("DataWriter is not *SQLiteWrapper; Telegram linking will be unavailable")
	}

	// Start the HTTP API server in the background (if enabled).
	if config.API.Enabled {
		emailSvc := util.NewEmailService(config.SMTP)
		if emailSvc.Enabled() {
			logger.Info("SMTP configured",
				zap.String("host", config.SMTP.Host),
				zap.Int("port", config.SMTP.Port),
				zap.String("from", config.SMTP.From),
				zap.Bool("auth", config.SMTP.Username != ""),
			)
		} else {
			logger.Warn("SMTP not configured — verification/reset links will be logged only (set SMTP_HOST to enable email)")
		}
		srv, err := api.NewServer(config, k.DataWriter, ozbscraper, cccscraper, logger, staticFiles, emailSvc)
		if err != nil {
			logger.Fatal("Failed to create API server", zap.Error(err))
		}

		go func() {
			if err := srv.Start(); err != nil {
				logger.Error("API server stopped", zap.Error(err))
			}
		}()

		if emailSvc.Enabled() && k.WebUserDB != nil {
			loc, err := time.LoadLocation(config.API.SummaryTimezone)
			if err != nil {
				logger.Warn("Invalid SUMMARY_TIMEZONE, falling back to UTC",
					zap.String("timezone", config.API.SummaryTimezone), zap.Error(err))
				loc = time.UTC
			}
			startDailySummaryScheduler(logger, k.WebUserDB, ozbscraper, cccscraper, emailSvc, loc)
		}

		// Graceful shutdown on SIGINT / SIGTERM — exits the whole process.
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			logger.Info("Shutdown signal received, exiting…")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				logger.Error("API server forced shutdown", zap.Error(err))
			}
			os.Exit(0)
		}()
	} else {
		// No API server — still handle signals so Ctrl+C works.
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			logger.Info("Shutdown signal received, exiting…")
			os.Exit(0)
		}()
	}

	// Start the Telegram bot (blocks until process exits).
	k.StartBot()
}

func startDailySummaryScheduler(
	logger *zap.Logger,
	webUserDB persist.WebUserDBIF,
	ozbScraper *scrapers.OzBargainScraper,
	cccScraper *scrapers.CamCamCamScraper,
	emailSvc *util.EmailService,
	loc *time.Location,
) {
	go func() {
		for {
			now := time.Now().In(loc)
			next := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, loc)
			if !now.Before(next) {
				next = next.Add(24 * time.Hour)
			}
			logger.Info("Daily summary scheduled", zap.Time("next_run", next))
			time.Sleep(time.Until(next))
			sendDailySummaries(logger, webUserDB, ozbScraper, cccScraper, emailSvc)
		}
	}()
}

func sendDailySummaries(
	logger *zap.Logger,
	webUserDB persist.WebUserDBIF,
	ozbScraper *scrapers.OzBargainScraper,
	cccScraper *scrapers.CamCamCamScraper,
	emailSvc *util.EmailService,
) {
	users, err := webUserDB.GetAllVerifiedWebUsers()
	if err != nil {
		logger.Error("daily summary: failed to fetch users", zap.Error(err))
		return
	}

	const maxDealAge = 24 * time.Hour

	// Collect OZB_SUPER deals posted within the last 24 hours, sorted by upvotes descending.
	// GetDealAge re-parses PostedOn at call time, so it reflects current age not scrape-time age.
	var ozbDeals []models.OzBargainDeal
	for _, d := range ozbScraper.Deals {
		if d.DealType == int(scrapers.OZB_SUPER) && ozbScraper.GetDealAge(d.PostedOn) <= maxDealAge {
			ozbDeals = append(ozbDeals, d)
		}
	}
	sort.Slice(ozbDeals, func(i, j int) bool {
		vi, _ := strconv.Atoi(ozbDeals[i].Upvotes)
		vj, _ := strconv.Atoi(ozbDeals[j].Upvotes)
		return vi > vj
	})

	// Collect AMZ_DAILY deals published within the last 24 hours.
	// gofeed sets Published as RFC1123Z or RFC1123; try both.
	amzLayouts := []string{time.RFC1123Z, time.RFC1123}
	cutoff := time.Now().Add(-maxDealAge)
	var amzDeals []models.CamCamCamDeal
	for _, d := range cccScraper.Deals {
		if d.DealType != int(scrapers.AMZ_DAILY) {
			continue
		}
		withinDay := false
		for _, layout := range amzLayouts {
			if t, err := time.Parse(layout, d.Published); err == nil {
				withinDay = t.After(cutoff)
				break
			}
		}
		if withinDay {
			amzDeals = append(amzDeals, d)
		}
	}

	sent := 0
	for _, user := range users {
		if !user.EmailSummary {
			continue
		}
		if err := emailSvc.SendDailySummary(user.Email, ozbDeals, amzDeals); err != nil {
			logger.Error("daily summary: send failed", zap.String("email", user.Email), zap.Error(err))
		} else {
			sent++
		}
	}
	logger.Info("daily summary sent", zap.Int("recipients", sent))
}
