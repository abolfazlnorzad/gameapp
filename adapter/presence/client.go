package presence

import (
	"context"
	"gameapp/contract/goproto/presence"
	"gameapp/dto"
	"gameapp/pkg/protobufmapper"
	"gameapp/pkg/slice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
}

func New() Client {
	return Client{}
}

func (c Client) GetPresence(ctx context.Context, request dto.GetPresenceRequest) (dto.GetPresenceResponse, error) {

	// TODO - use richerror on all adapter methods

	// TODO - what's the best practice for reliable communication
	// retry for connection time out?!
	// TODO -  is it okay to create new connection for every method call?
	conn, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return dto.GetPresenceResponse{}, err
	}
	defer conn.Close()

	client := presence.NewPresenceServiceClient(conn)
	response, err := client.GetPresence(ctx, &presence.GetPresenceRequest{
		UserIds: slice.MapUintToUint64(request.UserIDs),
	})
	if err != nil {
		return dto.GetPresenceResponse{}, err
	}
	return protobufmapper.MapGetPresenceResponseFromProtobuf(response), nil
}
