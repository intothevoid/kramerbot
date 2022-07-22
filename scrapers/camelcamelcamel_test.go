package scrapers

import (
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

func TestCamCamCamScraper_IsTargetDropGreater(t *testing.T) {
	type fields struct {
		BaseUrl         []string
		Logger          *zap.Logger
		SID             ScraperID
		ScrapeInterval  int
		MaxDealsToStore int
		Deals           []models.CamCamCamDeal
	}
	type args struct {
		deal   *models.CamCamCamDeal
		target int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				BaseUrl:         []string{"http://www.test.com"},
				Logger:          nil,
				SID:             1,
				ScrapeInterval:  1,
				MaxDealsToStore: 5,
				Deals:           []models.CamCamCamDeal{},
			},
			args: args{
				deal: &models.CamCamCamDeal{
					Id:        "deal1",
					Title:     "Canon EF 100-400mm f/4.5-5.6L IS II USM Lens, White - down 5.00% ($159.40) to $3,028.60 from $3,188.00",
					Url:       "http://www.test.com",
					Published: "2 hours ago",
					Image:     "",
					DealType:  5,
				},
				target: 4,
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				BaseUrl:         []string{"http://www.test.com"},
				Logger:          nil,
				SID:             1,
				ScrapeInterval:  1,
				MaxDealsToStore: 5,
				Deals:           []models.CamCamCamDeal{},
			},
			args: args{
				deal: &models.CamCamCamDeal{
					Id:        "deal2",
					Title:     "Yamaha NS-555 Floorstanding Sp...ss Reflex System, Black (Each) - down 27.85% ($584.52) to $1,514.48 from $2,099.00",
					Url:       "http://www.test.com",
					Published: "2 hours ago",
					Image:     "",
					DealType:  5,
				},
				target: 25,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CamCamCamScraper{
				BaseUrl:         tt.fields.BaseUrl,
				Logger:          tt.fields.Logger,
				SID:             tt.fields.SID,
				ScrapeInterval:  tt.fields.ScrapeInterval,
				MaxDealsToStore: tt.fields.MaxDealsToStore,
				Deals:           tt.fields.Deals,
			}
			if got := s.IsTargetDropGreater(tt.args.deal, tt.args.target); got != tt.want {
				t.Errorf("CamCamCamScraper.IsTargetDropGreater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCamCamCamScraper_GetDealDropString(t *testing.T) {
	type fields struct {
		BaseUrl         []string
		Logger          *zap.Logger
		SID             ScraperID
		ScrapeInterval  int
		MaxDealsToStore int
		Deals           []models.CamCamCamDeal
	}
	type args struct {
		deal *models.CamCamCamDeal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				BaseUrl:         []string{"http://www.test.com"},
				Logger:          nil,
				SID:             1,
				ScrapeInterval:  1,
				MaxDealsToStore: 5,
				Deals:           []models.CamCamCamDeal{},
			},
			args: args{
				deal: &models.CamCamCamDeal{
					Id:        "deal1",
					Title:     "Canon EF 100-400mm f/4.5-5.6L IS II USM Lens, White - down 5.00% ($159.40) to $3,028.60 from $3,188.00",
					Url:       "http://www.test.com",
					Published: "2 hours ago",
					Image:     "",
					DealType:  5,
				},
			},
			want: "down 5.00% ($159.40) to $3,028.60 from $3,188.00",
		},
		{
			name: "test2",
			fields: fields{
				BaseUrl:         []string{"http://www.test.com"},
				Logger:          nil,
				SID:             1,
				ScrapeInterval:  1,
				MaxDealsToStore: 5,
				Deals:           []models.CamCamCamDeal{},
			},
			args: args{
				deal: &models.CamCamCamDeal{
					Id:        "deal2",
					Title:     "NETGEAR AC797-100AUS 4G Router AirCard - down 23.57% ($76.78) to $249.00 from $325.78",
					Url:       "http://www.test.com",
					Published: "2 hours ago",
					Image:     "",
					DealType:  5,
				},
			},
			want: "down 23.57% ($76.78) to $249.00 from $325.78",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CamCamCamScraper{
				BaseUrl:         tt.fields.BaseUrl,
				Logger:          tt.fields.Logger,
				SID:             tt.fields.SID,
				ScrapeInterval:  tt.fields.ScrapeInterval,
				MaxDealsToStore: tt.fields.MaxDealsToStore,
				Deals:           tt.fields.Deals,
			}
			if got := s.GetDealDropString(tt.args.deal); got != tt.want {
				t.Errorf("CamCamCamScraper.GetDealDropString() = %v, want %v", got, tt.want)
			}
		})
	}
}
