package scrapers_test

import (
	"fmt"
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap/zapcore"
)

// Basic integration test — scrapes the live site.
func TestScrape(t *testing.T) {
	logger := util.SetupLogger(zapcore.DebugLevel, false)

	ozbScraper := new(scrapers.OzBargainScraper)
	ozbScraper.Logger = logger
	ozbScraper.BaseUrl = "https://www.ozbargain.com.au/"
	ozbScraper.ScrapeInterval = 5
	ozbScraper.MaxDealsToStore = 100

	t.Log("Scraping URL " + ozbScraper.BaseUrl)
	if err := ozbScraper.Scrape(); err != nil {
		t.Error("Error scraping:", err)
	}

	if len(ozbScraper.Deals) == 0 {
		t.Fatal("No deals found")
	}
	t.Log("Found " + fmt.Sprintf("%d", len(ozbScraper.Deals)) + " deals")
}

// TestGetDealAge verifies the age parser doesn't crash.
func TestGetDealAge(t *testing.T) {
	s := new(scrapers.OzBargainScraper)
	s.Logger = util.SetupLogger(zapcore.DebugLevel, false)
	t.Logf("Deal age: %s", s.GetDealAge("Neoika on 15/05/2022 - 14:38  kogan.com"))
}

// TestGetDealType_TopDeal asserts 25+ votes within 24h → OZB_SUPER.
func TestGetDealType_TopDeal(t *testing.T) {
	s := &scrapers.OzBargainScraper{Logger: util.SetupLogger(zapcore.DebugLevel, false)}

	deal := models.OzBargainDeal{
		Upvotes: "30",
		DealAge: "2h0m0s",
	}
	got := s.GetDealType(deal)
	if got != int(scrapers.OZB_SUPER) {
		t.Errorf("expected OZB_SUPER (%d) for 30 votes / 2h, got %d", scrapers.OZB_SUPER, got)
	}
}

// TestGetDealType_ExactThreshold asserts exactly 25 votes within 24h → OZB_SUPER.
func TestGetDealType_ExactThreshold(t *testing.T) {
	s := &scrapers.OzBargainScraper{Logger: util.SetupLogger(zapcore.DebugLevel, false)}

	deal := models.OzBargainDeal{
		Upvotes: "25",
		DealAge: "23h59m0s",
	}
	got := s.GetDealType(deal)
	if got != int(scrapers.OZB_SUPER) {
		t.Errorf("expected OZB_SUPER (%d) for 25 votes / 23h59m, got %d", scrapers.OZB_SUPER, got)
	}
}

// TestGetDealType_BelowThreshold asserts < 25 votes → OZB_REG.
func TestGetDealType_BelowThreshold(t *testing.T) {
	s := &scrapers.OzBargainScraper{Logger: util.SetupLogger(zapcore.DebugLevel, false)}

	deal := models.OzBargainDeal{
		Upvotes: "5",
		DealAge: "1h0m0s",
	}
	got := s.GetDealType(deal)
	if got != int(scrapers.OZB_REG) {
		t.Errorf("expected OZB_REG (%d) for 5 votes / 1h, got %d", scrapers.OZB_REG, got)
	}
}

// TestGetDealType_OldDeal asserts 100+ votes but older than 24h → OZB_REG.
func TestGetDealType_OldDeal(t *testing.T) {
	s := &scrapers.OzBargainScraper{Logger: util.SetupLogger(zapcore.DebugLevel, false)}

	deal := models.OzBargainDeal{
		Upvotes: "200",
		DealAge: "48h0m0s",
	}
	got := s.GetDealType(deal)
	if got != int(scrapers.OZB_REG) {
		t.Errorf("expected OZB_REG (%d) for 200 votes / 48h old, got %d", scrapers.OZB_REG, got)
	}
}

// TestFilter verifies keyword filtering works.
func TestFilter(t *testing.T) {
	s := &scrapers.OzBargainScraper{Logger: util.SetupLogger(zapcore.DebugLevel, false)}

	s.Deals = []models.OzBargainDeal{
		{Title: "Test deal", Upvotes: "49", DealAge: "0h59m00s"},
		{Title: "Test Beer Deal Weihenstephaner Schooner Cheap!", Upvotes: "100", DealAge: "5h59m00s"},
	}

	filtered := s.FilterByKeywords([]string{"w00t", "beer"})
	if len(filtered) == 0 {
		t.Fatal("Filter returned no deals — expected 1 match")
	}
	for _, d := range filtered {
		t.Log(d.Title)
	}
}

func TestOzBargainScraper_Scrape(t *testing.T) {
	logger := util.SetupLogger(zapcore.DebugLevel, false)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Test Scrape", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &scrapers.OzBargainScraper{
				Logger:          logger,
				BaseUrl:         "https://www.ozbargain.com.au/",
				ScrapeInterval:  5,
				MaxDealsToStore: 100,
			}
			if err := o.Scrape(); (err != nil) != tt.wantErr {
				t.Errorf("OzBargainScraper.Scrape() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
