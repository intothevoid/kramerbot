package util

import "testing"

func TestRssParser_ParseFeed(t *testing.T) {
	testUrl := "https://au.camelcamelcamel.com/top_drops/feed?t=daily&"

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
				Url: tt.fields.Url,
			}
			rss.ParseFeed()
		})
	}
}
