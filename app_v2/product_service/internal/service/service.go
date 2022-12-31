package service

import (
	"context"

	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/config"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/database/store"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Service struct {
	cfg     *config.Config
	log     logr.Logger
	storeDb store.StoreQuerier
	//api.UnimplementedOrderServiceServer
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
}

// RegisterWithHandler implementing service server interface
func (s *Service) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := api.NewHttpHandler(s)
	if err != nil {
		return err
	}
	return nil
}
