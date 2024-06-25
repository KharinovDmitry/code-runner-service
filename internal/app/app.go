package app

import (
	"code-runner-service/config"
	grpcServer "code-runner-service/internal/server/grpc"
	"code-runner-service/internal/service"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

func Run(cfg *config.Config) error {
	services := service.NewServiceManager()
	if err := services.Init(cfg.Env); err != nil {
		return err
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			services.Logger.Error(err.Error())
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
	))

	grpcServer.Register(server, services.TestRunner, services.Logger)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to open tcp conn: %w", err)
	}
	if err := server.Serve(l); err != nil {
		return fmt.Errorf("failed to run gprc server: %w", err)
	}

	services.Logger.Info("app started")
	return nil
}
