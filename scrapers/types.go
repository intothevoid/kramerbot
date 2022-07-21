package scrapers

// Deal types
type DealType int

const (
	OZB_REG DealType = iota
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
