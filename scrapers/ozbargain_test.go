package scrapers_test

import (
	"fmt"
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
)

// Basic test to check that the scraper works
func TestScrape(t *testing.T) {
	// create a new scraper
	ozbScraper := new(scrapers.OzBargainScraper)

	// initialise logger
	t.Log("Initialising logger")
	ozbScraper.Logger = util.SetupLogger()
	ozbScraper.BaseUrl = "https://www.ozbargain.com.au/"

	// Start scraping
	t.Log("Scraping URL " + ozbScraper.BaseUrl)
	ozbScraper.Scrape()

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

	inputstr := "Neoika on 15/05/2022 - 14:38  kogan.com"

	t.Logf("Deal age calculated was - %s", s.GetDealAge(inputstr))
}

// Test if deal is a super deal
func TestIsSuperDeal(t *testing.T) {
	// create a new scraper
	s := new(scrapers.OzBargainScraper)

	// create a new deal
	deal1 := models.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "49",
		DealAge:  "0h59m00s",
	}

	if s.GetDealType(deal1) == int(scrapers.SUPER_DEAL) {
		t.Log("Deal1 is a super deal")
	} else if s.GetDealType(deal1) == int(scrapers.GOOD_DEAL) {
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

	if s.GetDealType(deal2) == int(scrapers.SUPER_DEAL) {
		t.Log("Deal2 is a super deal")
	} else if s.GetDealType(deal2) == int(scrapers.GOOD_DEAL) {
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

	if s.GetDealType(deal3) == int(scrapers.SUPER_DEAL) {
		t.Log("Deal3 is a super deal")
	} else if s.GetDealType(deal2) == int(scrapers.GOOD_DEAL) {
		t.Log("Deal3 is a good deal")
	} else {
		t.Log("Deal3 is a regular deal")
	}
}

// Test if filter is working
func TestFilter(t *testing.T) {
	// create a new scraper
	s := new(scrapers.OzBargainScraper)

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
