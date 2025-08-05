package server

import (
	"context"
	"fmt"

	interceptor "github.com/kkiling/goplatform/interseptors"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"
	"google.golang.org/grpc"
)

// CustomHandlerService хендлер касромного сервер
type CustomHandlerService interface {
	server.HandlerService
}

// CustomServer касромный сервер
type CustomServer struct {
	server   *server.Server
	logger   log.Logger
	services []CustomHandlerService
}

// NewCustomServer новый сервер
func NewCustomServer(
	logger log.Logger,
	serverConfig server.Config,
	services ...CustomHandlerService,
) *CustomServer {
	return &CustomServer{
		server:   server.NewServer(logger, serverConfig),
		logger:   logger.Named("user_account_server"),
		services: services,
	}
}

func (p *CustomServer) unaryServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	return []grpc.UnaryServerInterceptor{
		interceptor.NewPanicRecoverInterceptor(p.logger),
		interceptor.NewLoggerInterceptor(p.logger),
	}, nil
}

// Start запустить сервер
func (p *CustomServer) Start(ctx context.Context, swaggerName string) error {
	unaryServerInterceptors, err := p.unaryServerInterceptors()
	if err != nil {
		return fmt.Errorf("p.unaryServerInterceptors(): %w", err)
	}

	p.server.WitUnaryServerInterceptor(unaryServerInterceptors...)

	var impl []server.HandlerService
	for _, service := range p.services {
		impl = append(impl, service)
	}

	if err = p.server.Start(ctx, swaggerName, impl...); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

// Stop остановить сервер
func (p *CustomServer) Stop() {
	p.server.Stop()
}
