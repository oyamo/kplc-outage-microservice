package server

import (
	"github.com/oyamo/kplc-outage-microservice/proto/notifications"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription/delivery"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	Config Config
	SubUC  subscription.UseCase
}

func NewServer(cfg Config, subc subscription.UseCase) *Server {
	return &Server{cfg, subc}
}

func (s *Server) Run(errChan chan<- error, exitChan chan<- int) {
	lis, err := net.Listen("tcp", ":"+s.Config.GRPCPort)
	if err != nil {
		errChan <- err
		exitChan <- -1
		return
	}

	// create grpc server
	grpcServer := grpc.NewServer()

	// create handler
	handler := delivery.NewHandler(s.SubUC)
	// Register handler
	notifications.RegisterNotificationsServer(grpcServer, handler)
	log.Printf("Serving GRPC on PORT %s", s.Config.GRPCPort)
	// Create server
	err = grpcServer.Serve(lis)
	if err != nil {
		errChan <- err
		exitChan <- -1
	}

	// Exit gracefully
	exitChan <- 0
}
