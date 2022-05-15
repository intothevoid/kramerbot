package scrapers_test

import (
	"fmt"
	"testing"

	"github.com/intothevoid/ozbot/scrapers"
	"github.com/intothevoid/ozbot/util"
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

	// create a new deal2
	deal1 := scrapers.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "49",
		DealAge:  "0h59m00s",
	}

	if s.IsSuperDeal(deal1) {
		t.Log("Deal1 is a super deal")
	} else {
		t.Log("Deal1 is not a super deal")
	}

	// create a new deal1
	deal2 := scrapers.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "55",
		DealAge:  "0h59m00s",
	}

	if s.IsSuperDeal(deal2) {
		t.Log("Deal2 is a super deal")
	} else {
		t.Log("Deal2 is not a super deal")
	}
}

// Test if deal is a good deal
func TestIsGoodDeal(t *testing.T) {
	// create a new scraper
	s := new(scrapers.OzBargainScraper)

	// create a new deal2
	deal1 := scrapers.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "24",
		DealAge:  "0h59m00s",
	}

	if s.IsGoodDeal(deal1) {
		t.Log("Deal1 is a good deal")
	} else {
		t.Log("Deal1 is not a good deal")
	}

	// create a new deal1
	deal2 := scrapers.OzBargainDeal{
		Title:    "Test deal",
		Url:      "https://www.ozbargain.com.au/deals/test-deal",
		PostedOn: "Neoika on 15/05/2022 - 14:38  kogan.com",
		Upvotes:  "30",
		DealAge:  "0h59m00s",
	}

	if s.IsGoodDeal(deal2) {
		t.Log("Deal2 is a good deal")
	} else {
		t.Log("Deal2 is not a good deal")
	}
}
