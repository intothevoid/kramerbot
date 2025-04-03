package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	persist "github.com/intothevoid/kramerbot/persist"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

func TestKramerBot_verifyAdminPassword(t *testing.T) {
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
	type args struct {
		message string
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
				Token:      "test",
				Logger:     nil,
				BotApi:     nil,
				OzbScraper: nil,
				CCCScraper: nil,
				UserStore:  nil,
				DataWriter: nil,
				Pipup:      nil,
				Config:     nil,
			},
			args: args{
				message: "testpassword:this is a test announcement",
			},
			want: false,
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
			if got := k.verifyAdminPassword(tt.args.message); got != tt.want {
				t.Errorf("KramerBot.verifyAdminPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
