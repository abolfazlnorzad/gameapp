package presence

import (
	"context"
	"gameapp/contract/golang/presence"
	"gameapp/dto"
	"gameapp/pkg/protobuf"
	"gameapp/pkg/slice"
	"google.golang.org/grpc"
)

type Client struct {
	client presence.PresenceServiceClient
}

func New(conn *grpc.ClientConn) Client {
	return Client{
		client: presence.NewPresenceServiceClient(conn),
	}
}

func (c Client) GetPresence(ctx context.Context, request dto.GetPresenceRequest) (dto.GetPresenceResponse, error) {
	response, err := c.client.GetPresence(ctx, &presence.GetPresenceRequest{
		UserIds: slice.MapUintToUint64(request.UserIDs),
	})
	if err != nil {
		return dto.GetPresenceResponse{}, err
	}
	return protobuf.MapGetPresenceResponseFromProtobuf(response), nil
}
