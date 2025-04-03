package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

type mockDatabase struct {
	persist.DatabaseIF
}

func (m *mockDatabase) WriteUserStore(userStore *models.UserStore) error {
	return nil
}

func (m *mockDatabase) ReadUserStore() (*models.UserStore, error) {
	return &models.UserStore{}, nil
}

func (m *mockDatabase) GetUser(chatID int64) (*models.UserData, error) {
	return &models.UserData{}, nil
}

func (m *mockDatabase) DeleteUser(user *models.UserData) error {
	return nil
}

func (m *mockDatabase) AddUser(user *models.UserData) error {
	return nil
}

func (m *mockDatabase) UpdateUser(user *models.UserData) error {
	return nil
}

func (m *mockDatabase) Close() error {
	return nil
}

func (m *mockDatabase) Ping() error {
	return nil
}

func TestKramerBot_processOzbargainDeals(t *testing.T) {
	type fields struct {
		Token      string
		Logger     *zap.Logger
		BotApi     *tgbotapi.BotAPI
		OzbScraper *scrapers.OzBargainScraper
		CCCScraper *scrapers.CamCamCamScraper
		UserStore  *models.UserStore
		DataWriter persist.DatabaseIF
		Pipup      *pipup.Pipup
		Config     *util.Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil scraper",
			fields: fields{
				Logger:     zap.NewNop(),
				OzbScraper: nil,
				UserStore:  &models.UserStore{Users: make(map[int64]*models.UserData)},
			},
			wantErr: true,
		},
		{
			name: "successful processing",
			fields: fields{
				Logger: zap.NewNop(),
				OzbScraper: &scrapers.OzBargainScraper{
					Logger:          zap.NewNop(),
					ScrapeInterval:  10,
					MaxDealsToStore: 50,
					BaseUrl:         "https://www.ozbargain.com.au",
					SID:             scrapers.SID_OZBARGAIN,
					Deals: []models.OzBargainDeal{
						{
							Id:       "123",
							Title:    "Test Deal 1",
							Url:      "https://example.com",
							Upvotes:  "30",
							DealType: int(scrapers.OZB_GOOD),
						},
						{
							Id:       "456",
							Title:    "Test Deal 2",
							Url:      "https://example.com",
							Upvotes:  "20",
							DealType: int(scrapers.OZB_GOOD),
						},
						{
							Id:       "789",
							Title:    "Nintendo Switch Deal",
							Url:      "https://example.com",
							Upvotes:  "20",
							DealType: int(scrapers.OZB_GOOD),
						},
					},
				},
				UserStore: &models.UserStore{
					Users: map[int64]*models.UserData{
						123: {
							ChatID:   123,
							Username: "testuser",
							OzbGood:  true,
							Keywords: []string{"test", "nintendo"},
							OzbSent:  []string{}, // Empty sent deals to verify multiple matches
						},
						456: {
							ChatID:   456,
							Username: "testuser2",
							OzbGood:  false,
							Keywords: []string{"switch"},
						},
					},
				},
				DataWriter: &mockDatabase{},
				Config: &util.Config{
					Scrapers: util.ScrapersConfig{
						OzBargain: util.OzBargainConfig{
							ScrapeInterval: 10,
							MaxStoredDeals: 50,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty user store",
			fields: fields{
				Logger: zap.NewNop(),
				OzbScraper: &scrapers.OzBargainScraper{
					Logger:          zap.NewNop(),
					ScrapeInterval:  10,
					MaxDealsToStore: 50,
					BaseUrl:         "https://www.ozbargain.com.au",
					SID:             scrapers.SID_OZBARGAIN,
					Deals: []models.OzBargainDeal{
						{
							Id:       "123",
							Title:    "Test Deal",
							Url:      "https://example.com",
							Upvotes:  "30",
							DealType: int(scrapers.OZB_GOOD),
						},
					},
				},
				UserStore:  &models.UserStore{Users: make(map[int64]*models.UserData)},
				DataWriter: &mockDatabase{},
				Config: &util.Config{
					Scrapers: util.ScrapersConfig{
						OzBargain: util.OzBargainConfig{
							ScrapeInterval: 10,
							MaxStoredDeals: 50,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KramerBot{
				Token:      tt.fields.Token,
				Logger:     tt.fields.Logger,
				BotApi:     tt.fields.BotApi,
				OzbScraper: tt.fields.OzbScraper,
				CCCScraper: tt.fields.CCCScraper,
				UserStore:  tt.fields.UserStore,
				DataWriter: tt.fields.DataWriter,
				Pipup:      tt.fields.Pipup,
				Config:     tt.fields.Config,
			}
			if err := k.processOzbargainDeals(); (err != nil) != tt.wantErr {
				t.Errorf("KramerBot.processOzbargainDeals() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
