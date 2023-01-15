package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/internal/database/store"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/internal/producer"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/util"
	"github.com/tabbed/pqtype"
	"google.golang.org/grpc/codes"
)

func (s *Service) AddProduct(ctx context.Context, req *api.AddProductRequest) (*api.AddProductResponse, error) {
	logger := s.log.WithName("AddProduct").WithValues("request", req)
	var variants pqtype.NullRawMessage
	err := json.Unmarshal(req.Variants, &variants)
	if err != nil {
		logger.Error(err, "Marshal json fail")
		return nil, err
	}
	id, err := s.storeDb.CreateProduct(ctx, store.CreateProductParams{
		TemplateID:  sql.NullInt64{Int64: req.TemplateId, Valid: true},
		Name:        req.Name,
		OriginPrice: req.OriginPrice,
		SalePrice:   req.SalePrice,
		State:       util.ProductStateAvailable,
		Variants:    variants,
		CreateUid:   req.CreateUid,
	})
	if err != nil {
		logger.Error(err, "CreateProduct")
		return nil, err
	}
	// Publish to kafka
	go func() {
		err := s.producer.Publish(ctx, "", producer.ProducerEvent{
			Key:   fmt.Sprintf("product/%v", id),
			Value: producer.InsertDatabaseEventValue(id),
		})
		if err != nil {
			s.log.Error(err, "Publish message fail")
			return
		}
	}()

	return &api.AddProductResponse{
		Code:    uint32(codes.OK),
		Message: "OK",
		Data: &api.AddProductResponse_Data{
			Id: id,
		},
	}, nil
}

func (s *Service) UpdateProduct(ctx context.Context, req *api.UpdateProductRequest) (res *api.UpdateProductResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteProduct(ctx context.Context, req *api.DeleteProductRequest) (res *api.DeleteProductResponse, err error) {
	//TODO implement me
	panic("implement me")
}
