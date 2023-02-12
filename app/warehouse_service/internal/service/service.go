package service

import (
	accountApi "Server-for-Ecommerce/app/account_service/api"
	orderApi "Server-for-Ecommerce/app/order_service/api"
	"Server-for-Ecommerce/app/warehouse_service/api"
	"Server-for-Ecommerce/app/warehouse_service/config"
	"Server-for-Ecommerce/app/warehouse_service/internal/store"
	"context"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
)

type Service struct {
	cfg           *config.Config
	log           logr.Logger
	storeDb       store.StoreQuerier
	orderClient   orderApi.OrderServiceClient
	accountClient accountApi.AccountServiceClient
	api.UnimplementedWarehouseServiceServer
}

func NewService(
	cfg *config.Config,
	log logr.Logger,
	storeDb store.StoreQuerier,
	orderClient orderApi.OrderServiceClient,
	accountClient accountApi.AccountServiceClient,
) *Service {
	return &Service{
		cfg:           cfg,
		log:           log,
		storeDb:       storeDb,
		orderClient:   orderClient,
		accountClient: accountClient,
	}
}

func (s *Service) Close(ctx context.Context) {
	s.storeDb.Close()
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
