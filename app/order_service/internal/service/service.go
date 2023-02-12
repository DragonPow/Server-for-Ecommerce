package service

import (
	"Server-for-Ecommerce/app/order_service/api"
	"Server-for-Ecommerce/app/order_service/config"
	"Server-for-Ecommerce/app/order_service/internal/cache"
	"Server-for-Ecommerce/app/order_service/internal/store"
	"context"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
)

type Service struct {
	cfg     *config.Config
	log     logr.Logger
	storeDb store.StoreQuerier
	api.UnimplementedOrderServiceServer
	cache cache.Cache
}

func NewService(
	cfg *config.Config,
	log logr.Logger,
	storeDb store.StoreQuerier,
	cache cache.Cache,
) *Service {
	return &Service{
		cfg:     cfg,
		log:     log,
		storeDb: storeDb,
		cache:   cache,
	}
}

func (s *Service) Close(ctx context.Context) {
	s.storeDb.Close()
}

// RegisterWithServer implementing service server interface
func (s *Service) RegisterWithServer(server *grpc.Server) {
	api.RegisterOrderServiceServer(server, s)
}

// RegisterWithHandler implementing service server interface
func (s *Service) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := api.RegisterOrderServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Println(err, "Error register servers")
	}

	return nil
}
