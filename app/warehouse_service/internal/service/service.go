package service

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/config"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
)

type Service struct {
	cfg *config.Config
	api.UnimplementedWarehouseServiceServer
}

func NewService(
	cfg *config.Config,
) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) Close(ctx context.Context) {
}

// RegisterWithServer implementing service server interface
func (s *Service) RegisterWithServer(server *grpc.Server) {
	api.RegisterWarehouseServiceServer(server, s)
}

// RegisterWithHandler implementing service server interface
func (s *Service) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := api.RegisterWarehouseServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Println(err, "Error register servers")
	}

	return nil
}
