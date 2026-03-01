package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/intothevoid/kramerbot/api/handlers"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
)

// buildOzbScraper returns a pre-seeded OzBargain scraper with one REG and one SUPER deal.
func buildOzbScraper() *scrapers.OzBargainScraper {
	s := &scrapers.OzBargainScraper{}
	s.Deals = []models.OzBargainDeal{
		{Id: "1", Title: "Regular Deal", Upvotes: "3", DealAge: "1h0m0s", DealType: int(scrapers.OZB_REG)},
		{Id: "2", Title: "Top Deal", Upvotes: "30", DealAge: "2h0m0s", DealType: int(scrapers.OZB_SUPER)},
	}
	return s
}

// buildAmazonScraper returns a pre-seeded CCC scraper with one daily and one weekly deal.
func buildAmazonScraper() *scrapers.CamCamCamScraper {
	s := &scrapers.CamCamCamScraper{}
	s.Deals = []models.CamCamCamDeal{
		{Id: "amz-1", Title: "Daily Drop", DealType: int(scrapers.AMZ_DAILY)},
		{Id: "amz-2", Title: "Weekly Drop", DealType: int(scrapers.AMZ_WEEKLY)},
	}
	return s
}

func getDealsFromResponse(t *testing.T, body []byte) []interface{} {
	t.Helper()
	var envelope struct {
		Data struct {
			Deals []interface{} `json:"deals"`
			Total int           `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("failed to parse response: %v\nbody: %s", err, body)
	}
	return envelope.Data.Deals
}

// TestGetOzbDeals_SuperFilter verifies that type=super returns only OZB_SUPER deals.
// This would have caught Bug 1 (DealType always OZB_REG in scraper cache).
func TestGetOzbDeals_SuperFilter(t *testing.T) {
	h := &handlers.Handler{OzbScraper: buildOzbScraper()}

	req := httptest.NewRequest(http.MethodGet, "/deals/ozbargain?type=super", nil)
	w := httptest.NewRecorder()
	h.GetOzbDeals(w, req)

	deals := getDealsFromResponse(t, w.Body.Bytes())
	if len(deals) != 1 {
		t.Errorf("type=super: expected 1 deal, got %d", len(deals))
	}
}

// TestGetOzbDeals_AllFilter verifies that no type param returns all deals.
func TestGetOzbDeals_AllFilter(t *testing.T) {
	h := &handlers.Handler{OzbScraper: buildOzbScraper()}

	req := httptest.NewRequest(http.MethodGet, "/deals/ozbargain", nil)
	w := httptest.NewRecorder()
	h.GetOzbDeals(w, req)

	deals := getDealsFromResponse(t, w.Body.Bytes())
	if len(deals) != 2 {
		t.Errorf("no filter: expected 2 deals, got %d", len(deals))
	}
}

// TestGetAmazonDeals_DailyFilter verifies type=daily returns only AMZ_DAILY deals.
func TestGetAmazonDeals_DailyFilter(t *testing.T) {
	h := &handlers.Handler{CCCScraper: buildAmazonScraper()}

	req := httptest.NewRequest(http.MethodGet, "/deals/amazon?type=daily", nil)
	w := httptest.NewRecorder()
	h.GetAmazonDeals(w, req)

	deals := getDealsFromResponse(t, w.Body.Bytes())
	if len(deals) != 1 {
		t.Errorf("type=daily: expected 1 deal, got %d", len(deals))
	}
}

// TestGetAmazonDeals_WeeklyFilter verifies type=weekly returns only AMZ_WEEKLY deals.
// This would catch Bug 2 if daily deals leaked into the weekly response.
func TestGetAmazonDeals_WeeklyFilter(t *testing.T) {
	h := &handlers.Handler{CCCScraper: buildAmazonScraper()}

	req := httptest.NewRequest(http.MethodGet, "/deals/amazon?type=weekly", nil)
	w := httptest.NewRecorder()
	h.GetAmazonDeals(w, req)

	deals := getDealsFromResponse(t, w.Body.Bytes())
	if len(deals) != 1 {
		t.Errorf("type=weekly: expected 1 deal, got %d", len(deals))
	}
}
