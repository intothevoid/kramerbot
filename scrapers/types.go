package scrapers

import "go.uber.org/zap"

// Deal types
type DealType int

var REGULAR_DEAL DealType = 0
var SUPER_DEAL DealType = 1
var GOOD_DEAL DealType = 2

// Scraper type
type ScraperID int

var SID_OZBARGAIN ScraperID = 0

type Scraper interface {
	Scrape()
	AutoScrape()
	GetData() interface{}
}

// Ozbargain scraper
type OzBargainScraper struct {
	BaseUrl        string
	Logger         *zap.Logger
	Deals          []OzBargainDeal
	SID            ScraperID // Scraper ID
	ScrapeInterval int       // Scrape interval
}

// Deal type
type OzBargainDeal struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	PostedOn string `json:"time"`
	Upvotes  string `json:"upvotes"`
	DealAge  string `json:"dealage"`
	DealType int    `json:"dealtype"`
}
