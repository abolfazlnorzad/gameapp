package presenceserver

import (
	"context"
	"fmt"
	"gameapp/contract/golang/presence"
	"gameapp/dto"
	"gameapp/pkg/protobuf"
	"gameapp/pkg/slice"
	"gameapp/service/precenseservice"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	presenceSvc precenseservice.Service
}

func New(pSvc precenseservice.Service) Server {
	return Server{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		presenceSvc:                        pSvc,
	}
}
func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	response, err := s.presenceSvc.GetPresence(ctx, dto.GetPresenceRequest{
		UserIDs: slice.MapUint64ToUint(req.UserIds),
	})
	if err != nil {
		return nil, err
	}
	return protobuf.MapGetPresenceResponseToProtobuf(response), nil
}

func (s Server) Start() {
	// generate new listener
	// todo => add port to config
	address := fmt.Sprintf(":%d", 8888)
	l, lErr := net.Listen("tcp", address)
	if lErr != nil {
		panic(lErr)
	}
	// generate new grpc server
	grpcServer := grpc.NewServer()
	// generate new pb presence server
	presenceSvcServer := Server{}
	// register pb presence server into grpc server
	presence.RegisterPresenceServiceServer(grpcServer, presenceSvcServer)

	log.Println("presence grpc server starting on", address)

	gErr := grpcServer.Serve(l)
	if gErr != nil {
		log.Fatal("couldn't server presence grpc server")
	}

}
