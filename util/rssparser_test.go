package util

import (
	"testing"

	"go.uber.org/zap"
)

func TestRssParser_ParseFeed(t *testing.T) {
	testUrl := "https://au.camelcamelcamel.com/top_drops/feed?t=daily&"

	// Create a test logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	type fields struct {
		Url string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{name: "test1", fields: fields{Url: testUrl}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rss := &RssParser{
				Url:    tt.fields.Url,
				Logger: logger,
			}
			feed, err := rss.ParseFeed()
			if err != nil {
				t.Errorf("ParseFeed() error = %v", err)
				return
			}
			if feed == nil {
				t.Error("ParseFeed() returned nil feed")
			}
		})
	}
}
