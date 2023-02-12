package service

import (
	"Server-for-Ecommerce/app_v2/db_manager_service/api"
	"Server-for-Ecommerce/app_v2/db_manager_service/internal/database/store"
	"Server-for-Ecommerce/app_v2/db_manager_service/producer"
	"Server-for-Ecommerce/app_v2/db_manager_service/util"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/tabbed/pqtype"
	"google.golang.org/grpc/codes"
	"time"
)

func (s *Service) AddProduct(ctx context.Context, req *api.AddProductRequest) (*api.AddProductResponse, error) {
	logger := s.log.WithName("AddProduct").WithValues("request", req)
	var variants pqtype.NullRawMessage
	variants.RawMessage = []byte(req.Variants)
	if req.Variants != util.EmptyString {
		variants.Valid = true
	}
	id, err := s.storeDb.CreateProduct(ctx, store.CreateProductParams{
		TemplateID:  sql.NullInt64{Int64: req.TemplateId, Valid: true},
		Image:       req.Image,
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
	go func(logger logr.Logger) {
		ctx := context.Background()
		err := s.producer.Publish(ctx, util.TopicInsertProduct, producer.ProducerEvent{
			Key:   fmt.Sprintf("product/%v", id),
			Value: producer.InsertDatabaseEventValue(id),
		})
		if err != nil {
			logger.Error(err, "Publish message fail")
			return
		}
	}(logger)

	return &api.AddProductResponse{
		Code:    uint32(codes.OK),
		Message: "OK",
		Data: &api.AddProductResponse_Data{
			Id: id,
		},
	}, nil
}

func (s *Service) UpdateProduct(ctx context.Context, req *api.UpdateProductRequest) (res *api.UpdateProductResponse, err error) {
	logger := s.log.WithName("UpdateProduct").WithValues("request", req)

	reader := bytes.NewReader(req.Variants)
	decoder := json.NewDecoder(reader)
	var updateRequestParams UpdateProductParams
	err = decoder.Decode(&updateRequestParams)
	if err != nil {
		logger.Error(err, "Decode variants fail")
		return nil, err
	}

	updateRequestParams.ID = req.Id
	writeTime, err := s.storeDb.UpdateProduct(ctx, updateRequestParams.ToStore())
	if err != nil {
		logger.Error(err, "UpdateProduct")
		return nil, err
	}
	// Publish to kafka
	go func(logger logr.Logger, writeTime time.Time, id int64, variants []byte) {
		ctx := context.Background()
		err = s.producer.Publish(ctx, util.TopicUpdateProduct, producer.ProducerEvent{
			Key: fmt.Sprintf("product/%v", req.Id),
			Value: producer.UpdateDatabaseEventValue{
				Id:         id,
				Variants:   variants,
				TimeUpdate: writeTime,
			},
		})
		if err != nil {
			logger.Error(err, "Publish message fail")
			return
		}
	}(logger, writeTime, req.Id, req.Variants)

	return &api.UpdateProductResponse{
		Code:    uint32(codes.OK),
		Message: "OK",
	}, nil
}

func (s *Service) DeleteProduct(ctx context.Context, req *api.DeleteProductRequest) (res *api.DeleteProductResponse, err error) {
	//TODO implement me
	panic("implement me")
}

type UpdateProductParams struct {
	TemplateID  util.NullInt64   `json:"template_id,omitempty"`
	Name        util.NullString  `json:"name,omitempty"`
	OriginPrice util.NullFloat64 `json:"origin_price,omitempty"`
	SalePrice   util.NullFloat64 `json:"sale_price,omitempty"`
	State       util.NullString  `json:"state,omitempty"`
	Variants    map[string]any   `json:"variants,omitempty"`
	CreateUid   int64            `json:"create_uid,omitempty"`
	ID          int64            `json:"id"`
}

func (u *UpdateProductParams) ToStore() store.UpdateProductParams {
	b, _ := json.Marshal(u.Variants)
	variants := pqtype.NullRawMessage{
		RawMessage: b,
		Valid:      len(u.Variants) > util.ZeroLength,
	}
	return store.UpdateProductParams{
		TemplateID:  u.TemplateID.NullInt64,
		Name:        u.Name.NullString,
		OriginPrice: u.OriginPrice.NullFloat64,
		SalePrice:   u.SalePrice.NullFloat64,
		State:       u.State.NullString,
		Variants:    variants,
		CreateUid:   u.CreateUid,
		ID:          u.ID,
	}
}
