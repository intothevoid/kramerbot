package handlers

import (
	"net/http"
	"strconv"

	"github.com/intothevoid/kramerbot/scrapers"
)

// GetOzbDeals returns OzBargain deals from the scraper's in-memory cache.
// Query params: type=good|super|all (default: all), limit (default 50), offset (default 0).
func (h *Handler) GetOzbDeals(w http.ResponseWriter, r *http.Request) {
	dealType := r.URL.Query().Get("type")
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)

	if h.OzbScraper == nil {
		jsonOK(w, map[string]interface{}{"deals": []interface{}{}, "total": 0})
		return
	}

	var filtered []interface{}
	for i := range h.OzbScraper.Deals {
		d := h.OzbScraper.Deals[i]
		switch dealType {
		case "good":
			if d.DealType == int(scrapers.OZB_GOOD) {
				filtered = append(filtered, d)
			}
		case "super":
			if d.DealType == int(scrapers.OZB_SUPER) {
				filtered = append(filtered, d)
			}
		default:
			filtered = append(filtered, d)
		}
	}

	total := len(filtered)
	filtered = paginate(filtered, offset, limit)
	jsonOK(w, map[string]interface{}{"deals": filtered, "total": total})
}

// GetAmazonDeals returns Amazon deals from the scraper's in-memory cache.
// Query params: type=daily|weekly|all (default: all), limit, offset.
func (h *Handler) GetAmazonDeals(w http.ResponseWriter, r *http.Request) {
	dealType := r.URL.Query().Get("type")
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)

	if h.CCCScraper == nil {
		jsonOK(w, map[string]interface{}{"deals": []interface{}{}, "total": 0})
		return
	}

	var filtered []interface{}
	for i := range h.CCCScraper.Deals {
		d := h.CCCScraper.Deals[i]
		switch dealType {
		case "daily":
			if d.DealType == int(scrapers.AMZ_DAILY) {
				filtered = append(filtered, d)
			}
		case "weekly":
			if d.DealType == int(scrapers.AMZ_WEEKLY) {
				filtered = append(filtered, d)
			}
		default:
			filtered = append(filtered, d)
		}
	}

	total := len(filtered)
	filtered = paginate(filtered, offset, limit)
	jsonOK(w, map[string]interface{}{"deals": filtered, "total": total})
}

// GetAllDeals returns a combined OzBargain + Amazon deal feed.
func (h *Handler) GetAllDeals(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 100)
	offset := queryInt(r, "offset", 0)

	var combined []interface{}
	if h.OzbScraper != nil {
		for i := range h.OzbScraper.Deals {
			combined = append(combined, h.OzbScraper.Deals[i])
		}
	}
	if h.CCCScraper != nil {
		for i := range h.CCCScraper.Deals {
			combined = append(combined, h.CCCScraper.Deals[i])
		}
	}

	total := len(combined)
	combined = paginate(combined, offset, limit)
	jsonOK(w, map[string]interface{}{"deals": combined, "total": total})
}

func queryInt(r *http.Request, param string, defaultVal int) int {
	s := r.URL.Query().Get(param)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return defaultVal
	}
	return v
}

func paginate(items []interface{}, offset, limit int) []interface{} {
	if offset >= len(items) {
		return []interface{}{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}
