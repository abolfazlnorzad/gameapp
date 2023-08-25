package presenceserver

import (
	"context"
	"fmt"
	"gameapp/adapter/redisadapter"
	"gameapp/config"
	"gameapp/contract/goproto/presence"
	"gameapp/dto"
	"gameapp/pkg/protobufmapper"
	"gameapp/pkg/slice"
	"gameapp/repository/redis/redispresence"
	"gameapp/service/presenceservice"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	psv presenceservice.Service
}

func New(pSvc presenceservice.Service) Server {
	return Server{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		psv:                                pSvc,
	}
}
func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	fmt.Println("req", req)
	cfg := config.Load()
	redisAdp := redisadapter.New(cfg.Redis)
	pr := redispresence.New(redisAdp)
	presenceSvc := presenceservice.New(pr, cfg.PresenceService)
	response, err := presenceSvc.GetPresence(ctx, dto.GetPresenceRequest{
		UserIDs: slice.MapUint64ToUint(req.UserIds),
	})
	if err != nil {
		return nil, err
	}
	return protobufmapper.MapGetPresenceResponseToProtobuf(response), nil
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
