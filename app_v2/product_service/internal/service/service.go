package service

import (
	"Server-for-Ecommerce/app_v2/product_service/cache/mem_cache"
	"Server-for-Ecommerce/app_v2/product_service/cache/redis_cache"
	"Server-for-Ecommerce/app_v2/product_service/database/store"
	"context"
	"net/http"
	"sync"

	"Server-for-Ecommerce/app_v2/product_service/api"
	"Server-for-Ecommerce/app_v2/product_service/config"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Service struct {
	cfg        *config.Config
	log        logr.Logger
	storeDb    store.StoreQuerier
	localCache redis_cache.RedisCache
	memCache   mem_cache.MemCache
	lockCache  LockCache
	//api.UnimplementedOrderServiceServer
}

func NewService(
	cfg *config.Config,
	log logr.Logger,
	storeDb store.StoreQuerier,
	cache redis_cache.RedisCache,
	memCache mem_cache.MemCache,
) *Service {
	return &Service{
		cfg:        cfg,
		log:        log,
		storeDb:    storeDb,
		localCache: cache,
		memCache:   memCache,
	}
}

func (s *Service) Close(ctx context.Context) {
	s.storeDb.Close()
	s.localCache.Close()
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

type LockCache struct {
	list sync.Map
	mu   sync.RWMutex
}
