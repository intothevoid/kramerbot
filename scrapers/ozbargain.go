package scrapers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

// Ozbargain scraper
// implement the scraper interface
type OzBargainScraper struct {
	BaseUrl string
	Logger  *zap.Logger
	Deals   []OzBargainDeal
}

type OzBargainDeal struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	PostedOn string `json:"time"`
	Upvotes  string `json:"upvotes"`
	DealAge  string `json:"dealage"`
	DealType int    `json:"dealtype"`
}

// Scrape the url
func (s *OzBargainScraper) Scrape() {
	url := s.BaseUrl + "/deals"
	s.Deals = []OzBargainDeal{}

	// create a new collector
	c := colly.NewCollector()

	// find the title class
	c.OnHTML("div .node.node-ozbdeal.node-teaser", func(e *colly.HTMLElement) {

		// get the formatted title and url
		dealTitle := e.ChildAttr(".n-right h2.title", "data-title")

		// get the deal url
		dealURL := s.BaseUrl[:len(s.BaseUrl)-1] + e.ChildAttr(".n-right h2 a", "href")

		// get the deal poster and time
		postedOn := e.ChildText(".n-right div.submitted")

		// get the deal upvotes
		upVotes := e.ChildText(".n-left .n-vote.n-deal.inact .nvb.voteup")

		// populate the deal
		deal := OzBargainDeal{
			Title:    dealTitle,
			Url:      dealURL,
			PostedOn: postedOn,
			Upvotes:  upVotes,
			DealAge:  s.GetDealAge(postedOn).String(),
			DealType: int(REGULAR_DEAL), //s.GetDealType(dealTitle, upVotes),
		}
		s.Logger.Debug("Found deal", zap.String("title", deal.Title), zap.String("url", deal.Url), zap.String("time", deal.PostedOn))

		// create item list
		s.Deals = append(s.Deals, deal)
	})

	// Start scraping
	c.Visit(url)
}

// Calculate the time elapsed since the deal was posted
func (s *OzBargainScraper) GetDealAge(postedOn string) time.Duration {
	// regular expression to pull time from string
	// re := regexp.MustCompile(`\d[\d\/:\s-]+\d`)
	re := regexp.MustCompile(`[\d\/]+\s*\-\s*[\d:]+`)
	dealTimestamp := re.FindString(postedOn)

	// time format as scraped from ozbargain
	const layout = "02/01/2006 - 15:04"

	tmts, err := time.Parse(layout, dealTimestamp)
	if err != nil {
		s.Logger.Error("Error parsing time", zap.Error(err))
	}

	tmnow, err := time.Parse(layout, time.Now().Format(layout))
	if err != nil {
		s.Logger.Error("Error parsing time", zap.Error(err))
	}

	// duration since post = time now - timestamp
	return tmnow.Sub(tmts)
}

// Check if deal is a super deal, good deal or just a regular deal
func (s *OzBargainScraper) GetDealType(deal OzBargainDeal) int {
	upvotes := deal.Upvotes
	dealAge := deal.DealAge

	duration, err := time.ParseDuration(dealAge)
	if err != nil {
		s.Logger.Error("Error parsing time", zap.Error(err))
	}

	// convert upvotes to int
	upvotesInt, err := strconv.Atoi(upvotes)

	if err != nil {
		s.Logger.Error("Error converting upvotes to int", zap.Error(err))
	}

	// 50+ upvotes within an hour
	if duration.Hours() < 1.0 && upvotesInt >= 50 {
		return int(SUPER_DEAL)
	}

	// 25+ upvotes within an hour
	if duration.Hours() < 1.0 && upvotesInt >= 25 {
		return int(GOOD_DEAL)
	}

	// regular deal
	return int(REGULAR_DEAL)
}

// Filter list of deals by keywords
func (s *OzBargainScraper) FilterByKeywords(keywords []string) []OzBargainDeal {
	filteredDeals := []OzBargainDeal{}
	for _, deal := range s.Deals {
		for _, keyword := range keywords {
			if strings.Contains(deal.Title, keyword) {
				filteredDeals = append(filteredDeals, deal)
			}
		}
	}
	return filteredDeals
}
