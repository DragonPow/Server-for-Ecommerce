package service

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/config"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/internal/store"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
)

type Service struct {
	cfg     *config.Config
	log     logr.Logger
	storeDb store.StoreQuerier
	api.UnimplementedWarehouseServiceServer
}

func NewService(
	cfg *config.Config,
	log logr.Logger,
	storeDb store.StoreQuerier,
) *Service {
	return &Service{
		cfg:     cfg,
		log:     log,
		storeDb: storeDb,
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
