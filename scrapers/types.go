package scrapers

// Deal types
type DealType int

const (
	UNKNOWN DealType = iota
	OZB_REG
	OZB_SUPER
	OZB_GOOD
	AMZ_DAILY
	AMZ_WEEKLY
)

// Scraper type
type ScraperID int

type Scraper interface {
	Scrape()
	AutoScrape()
	GetData() interface{}
}
