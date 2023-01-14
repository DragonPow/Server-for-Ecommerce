package service

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache/mem_cache"
	"net/http"

	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/config"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/database/store"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Service struct {
	cfg      *config.Config
	log      logr.Logger
	storeDb  store.StoreQuerier
	cache    cache.Cache
	memCache mem_cache.MemCache
	//api.UnimplementedOrderServiceServer
}

func NewService(
	cfg *config.Config,
	log logr.Logger,
	storeDb store.StoreQuerier,
	cache cache.Cache,
	memCache mem_cache.MemCache,
) *Service {
	return &Service{
		cfg:      cfg,
		log:      log,
		storeDb:  storeDb,
		cache:    cache,
		memCache: memCache,
	}
}

func (s *Service) Close(ctx context.Context) {
	s.storeDb.Close()
}

// RegisterWithServer implementing service server interface
func (s *Service) RegisterWithServer(server *grpc.Server) {
}

func (s *Service) RegisterWithHttpHandler(httpPattern string) (http.Handler, error) {
	return api.NewHttpHandler(httpPattern, s), nil
}

// RegisterWithHandler implementing service server interface
func (s *Service) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return nil
}
