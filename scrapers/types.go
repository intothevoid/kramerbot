package scrapers

// Deal types
type DealType int

var REGULAR_DEAL DealType = 0
var SUPER_DEAL DealType = 1
var GOOD_DEAL DealType = 2

// Scraper type
type ScraperID int

type Scraper interface {
	Scrape()
	AutoScrape()
	GetData() interface{}
}
