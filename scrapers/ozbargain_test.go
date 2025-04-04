package scrapers_test

import (
	"fmt"
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Basic test to check that the scraper works
func TestScrape(t *testing.T) {
	// Create a test logger
	logger := util.SetupLogger(zapcore.DebugLevel, false)

	// Create a test config
	config := &util.Config{
		LogLevel:  -1,
		LogToFile: false,
		TestMode:  true,
		Scrapers: util.ScrapersConfig{
			OzBargain: util.OzBargainConfig{
				ScrapeInterval: 5,
				MaxStoredDeals: 100,
			},
		},
	}

	// create a new scraper
	ozbScraper := new(scrapers.OzBargainScraper)
	ozbScraper.Logger = logger
	ozbScraper.BaseUrl = "https://www.ozbargain.com.au/"
	ozbScraper.ScrapeInterval = config.Scrapers.OzBargain.ScrapeInterval
	ozbScraper.MaxDealsToStore = config.Scrapers.OzBargain.MaxStoredDeals

	// Start scraping
	t.Log("Scraping URL " + ozbScraper.BaseUrl)
	err := ozbScraper.Scrape()
	if err != nil {
		t.Error("Error scraping.", err.Error())
	}

	if len(ozbScraper.Deals) == 0 {
		t.Error("No deals found")
		t.Fail()
	} else {
		t.Log("Found " + fmt.Sprintf("%d", len(ozbScraper.Deals)) + " deals")
		for deal := range ozbScraper.Deals {
			t.Log(ozbScraper.Deals[deal].Title)
		}
	}
}

// Test the getDealAge function
func TestGetDealAge(t *testing.T) {
	// create a new scraper
	s := new(scrapers.OzBargainScraper)
	s.Logger = util.SetupLogger(zapcore.DebugLevel, false)

	inputstr := "Neoika on 15/05/2022 - 14:38  kogan.com"

	t.Logf("Deal age calculated was - %s", s.GetDealAge(inputstr))
}

// Test if deal is a super deal
func TestIsSuperDeal(t *testing.T) {
	// create a new scraper
	s := new(scrapers.OzBargainScraper)
	s.Logger = util.SetupLogger(zapcore.DebugLevel, false)

	// create a new deal
	deal1 := models.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "49",
		DealAge:  "0h59m00s",
	}

	if s.GetDealType(deal1) == int(scrapers.OZB_SUPER) {
		t.Log("Deal1 is a super deal")
	} else if s.GetDealType(deal1) == int(scrapers.OZB_GOOD) {
		t.Log("Deal1 is a good deal")
	} else {
		t.Log("Deal1 is a regular deal")
	}

	// create a new deal
	deal2 := models.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "55",
		DealAge:  "0h59m00s",
	}

	if s.GetDealType(deal2) == int(scrapers.OZB_SUPER) {
		t.Log("Deal2 is a super deal")
	} else if s.GetDealType(deal2) == int(scrapers.OZB_GOOD) {
		t.Log("Deal2 is a good deal")
	} else {
		t.Log("Deal2 is a regular deal")
	}

	// create a new deal
	deal3 := models.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "20",
		DealAge:  "0h25m00s",
	}

	if s.GetDealType(deal3) == int(scrapers.OZB_SUPER) {
		t.Log("Deal3 is a super deal")
	} else if s.GetDealType(deal2) == int(scrapers.OZB_GOOD) {
		t.Log("Deal3 is a good deal")
	} else {
		t.Log("Deal3 is a regular deal")
	}
}

// Test if filter is working
func TestFilter(t *testing.T) {
	// create a new scraper
	s := new(scrapers.OzBargainScraper)
	s.Logger = util.SetupLogger(zapcore.DebugLevel, false)

	// create a new deal
	deal1 := models.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "49",
		DealAge:  "0h59m00s",
	}

	deal2 := models.OzBargainDeal{
		Title:    "Test Beer Deal Weihenstephaner Schooner Cheap!",
		Url:      "https://www.ozbargain.com.au/deals/test-deal-beer",
		PostedOn: "intothevoid on 15/05/2022 - 14:38  danmurphys.com",
		Upvotes:  "100",
		DealAge:  "5h59m00s",
	}

	s.Deals = []models.OzBargainDeal{deal1, deal2}

	filtered := s.FilterByKeywords([]string{"w00t", "beer"})

	if len(filtered) == 0 {
		t.Error("Filter deals did not work as intended")
		t.Fail()
	}

	for deal := range filtered {
		t.Log(filtered[deal].Title)
	}
}

func TestOzBargainScraper_Scrape(t *testing.T) {
	// Create a test logger
	logger := util.SetupLogger(zapcore.DebugLevel, false)

	// Create a test config
	config := &util.Config{
		LogLevel:  -1,
		LogToFile: false,
		TestMode:  true,
		Scrapers: util.ScrapersConfig{
			OzBargain: util.OzBargainConfig{
				ScrapeInterval: 5,
				MaxStoredDeals: 100,
			},
		},
	}

	type fields struct {
		Logger          *zap.Logger
		BaseUrl         string
		Deals           []models.OzBargainDeal
		SID             scrapers.ScraperID
		ScrapeInterval  int
		MaxDealsToStore int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test Scrape",
			fields: fields{
				Logger:          logger,
				BaseUrl:         "https://www.ozbargain.com.au/",
				ScrapeInterval:  config.Scrapers.OzBargain.ScrapeInterval,
				MaxDealsToStore: config.Scrapers.OzBargain.MaxStoredDeals,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &scrapers.OzBargainScraper{
				Logger:          tt.fields.Logger,
				BaseUrl:         tt.fields.BaseUrl,
				Deals:           tt.fields.Deals,
				SID:             tt.fields.SID,
				ScrapeInterval:  tt.fields.ScrapeInterval,
				MaxDealsToStore: tt.fields.MaxDealsToStore,
			}
			if err := o.Scrape(); (err != nil) != tt.wantErr {
				t.Errorf("OzBargainScraper.Scrape() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
