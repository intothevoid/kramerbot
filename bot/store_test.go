package bot

import (
	"testing"

	"github.com/intothevoid/kramerbot/models"
)

func TestOzbDealSent(t *testing.T) {
	type args struct {
		user *models.UserData
		deal *models.OzBargainDeal
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Deal already sent",
			args: args{
				user: &models.UserData{
					OzbSent: []string{"123"},
				},
				deal: &models.OzBargainDeal{
					Id: "123",
				},
			},
			want: true,
		},
		{
			name: "Deal not sent",
			args: args{
				user: &models.UserData{
					OzbSent: []string{"456"},
				},
				deal: &models.OzBargainDeal{
					Id: "123",
				},
			},
			want: false,
		},
		{
			name: "Empty sent list",
			args: args{
				user: &models.UserData{
					OzbSent: []string{},
				},
				deal: &models.OzBargainDeal{
					Id: "123",
				},
			},
			want: false,
		},
		{
			name: "Multiple deals sent",
			args: args{
				user: &models.UserData{
					OzbSent: []string{"123", "456", "789"},
				},
				deal: &models.OzBargainDeal{
					Id: "456",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OzbDealSent(tt.args.user, tt.args.deal); got != tt.want {
				t.Errorf("OzbDealSent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmzDealSent(t *testing.T) {
	type args struct {
		user *models.UserData
		deal *models.CamCamCamDeal
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Deal already sent",
			args: args{
				user: &models.UserData{
					AmzSent: []string{"123"},
				},
				deal: &models.CamCamCamDeal{
					Id: "123",
				},
			},
			want: true,
		},
		{
			name: "Deal not sent",
			args: args{
				user: &models.UserData{
					AmzSent: []string{"456"},
				},
				deal: &models.CamCamCamDeal{
					Id: "123",
				},
			},
			want: false,
		},
		{
			name: "Empty sent list",
			args: args{
				user: &models.UserData{
					AmzSent: []string{},
				},
				deal: &models.CamCamCamDeal{
					Id: "123",
				},
			},
			want: false,
		},
		{
			name: "Multiple deals sent",
			args: args{
				user: &models.UserData{
					AmzSent: []string{"123", "456", "789"},
				},
				deal: &models.CamCamCamDeal{
					Id: "456",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AmzDealSent(tt.args.user, tt.args.deal); got != tt.want {
				t.Errorf("AmzDealSent() = %v, want %v", got, tt.want)
			}
		})
	}
}
