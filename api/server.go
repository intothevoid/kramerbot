package api

import (
	"context"

	"google.golang.org/grpc"

	apiv1 "github.com/intothevoid/kramerbot/api/v1"
)

type KramerServiceServer struct {
}

func (s *KramerServiceServer) GetUserDetails(ctx context.Context, req *apiv1.UserDetailsRequest, opts ...grpc.CallOption) (*apiv1.UserDetailsResponse, error) {
	chatID := req.Chatid

	resp := apiv1.UserDetailsResponse{
		Chatid:     chatID,
		Username:   "test",
		Gooddeals:  true,
		Superdeals: true,
		Keywords:   "test keywords string",
		Dealssent:  "test deals sent string",
	}

	return &resp, nil
}
