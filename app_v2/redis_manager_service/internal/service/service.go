package service

import (
	"Server-for-Ecommerce/app_v2/redis_manager_service/internal/redis"
	producer "Server-for-Ecommerce/library/kafka/pub"
	"context"
	"github.com/gorilla/mux"
	"net/http"

	"Server-for-Ecommerce/app_v2/redis_manager_service/config"
	"Server-for-Ecommerce/app_v2/redis_manager_service/internal/database/store"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Service struct {
	cfg      *config.Config
	log      logr.Logger
	storeDb  store.StoreQuerier
	redis    *redis.Redis
	producer producer.Producer
	//api.UnimplementedOrderServiceServer
}

func NewService(
	cfg *config.Config,
	log logr.Logger,
	storeDb store.StoreQuerier,
	redis *redis.Redis,
	producer producer.Producer,
) *Service {
	return &Service{
		cfg:      cfg,
		log:      log,
		storeDb:  storeDb,
		redis:    redis,
		producer: producer,
	}
}

func (s *Service) Close(ctx context.Context) {
	s.storeDb.Close()
}

// RegisterWithServer implementing service server interface
func (s *Service) RegisterWithServer(server *grpc.Server) {
}

func (s *Service) RegisterWithHttpHandler(httpPattern string) (http.Handler, error) {
	return mux.NewRouter().PathPrefix(httpPattern).Subrouter(), nil
}

// RegisterWithHandler implementing service server interface
func (s *Service) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return nil
}
