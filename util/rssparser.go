package util

import (
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

type RssParser struct {
	Url    string
	Logger *zap.Logger
}

// Parse the RSS Url and return feed item
func (rss *RssParser) ParseFeed() (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss.Url)
	if err != nil {
		rss.Logger.Error("error parsing feed", zap.String("url", rss.Url), zap.Error(err))
		return nil, err
	}

	return feed, err
}

// Change Url
func (rss *RssParser) SetUrl(url string) {
	rss.Url = url
}
