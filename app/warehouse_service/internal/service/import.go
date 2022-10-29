package service

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/api"
	"google.golang.org/grpc/codes"
)

func (s *Service) CreateImportRequest(ctx context.Context, request *api.CreateImportRequestRequest) (*api.CreateImportRequestResponse, error) {
	return &api.CreateImportRequestResponse{
		Code:    int32(codes.OK),
		Message: request.String(),
		Data:    nil,
	}, nil
}
