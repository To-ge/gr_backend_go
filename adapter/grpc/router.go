package grpc

import (
	"fmt"
	"log"
	"net"

	v1 "github.com/To-ge/gr_backend_go/adapter/grpc/api/gen/go/v1"
	"github.com/To-ge/gr_backend_go/adapter/grpc/handler"
	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/repository"
	"github.com/To-ge/gr_backend_go/usecase"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func InitRouter() error {
	var err error

	address := config.LoadConfig().GrpcInfo.Address

	listenPort, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryHandler),
	}
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)

	dbConn, err := database.NewDBConnector()
	if err != nil {
		return err
	}

	v1.RegisterTelemetryServiceServer(server, handler.NewTelemetryHandler(usecase.NewTelemetryUsecase(service.NewTelemetryService(repository.NewtelemetryRepository(dbConn), repository.NewTelemetryLogRepository(dbConn)))))

	reflection.Register(server)
	go server.Serve(listenPort)
	log.Println("grpc server is running! addr: ", address)

	return nil
}

func recoveryHandler(p interface{}) error {
	log.Printf("Recovered from panic: %v", p)
	return status.Errorf(codes.Internal, "unexpected error occurred")
}
