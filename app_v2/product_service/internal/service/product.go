package service

import (
	"context"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetDetailProduct(ctx context.Context, req *api.GetDetailProductRequest) (res *api.GetDetailProductResponse, err error) {
	logger := s.log.WithName("GetDetailProduct").WithValues("request", req)
	// Get from memory
	memCacheProduct, ok := s.memCache.GetProduct(req.Id)
	if ok {
		// Get Template
		template, err := getOrInsertMem(s, ctx, memCacheProduct)
		if err != nil {
			logger.Error(err, "getOrInsertMem")
			return nil, err
		}

		return &api.GetDetailProductResponse{
			Code:    0,
			Message: "OK",
			Data: &api.ProductDetail{
				Id:                  memCacheProduct.ID,
				Name:                memCacheProduct.Name,
				OriginPrice:         memCacheProduct.OriginPrice,
				SalePrice:           memCacheProduct.SalePrice,
				Variants:            memCacheProduct.Variants,
				CreatedBy:           "",
				CreatedDate:         util.ParseTimeToString(memCacheProduct.CreateDate),
				UpdatedBy:           "",
				UpdatedDate:         util.ParseTimeToString(memCacheProduct.WriteDate),
				TemplateId:          memCacheProduct.TemplateID,
				TemplateName:        template.Name,
				TemplateDescription: template.Description,
				SoldQuantity:        template.SoldQuantity,
				RemainQuantity:      template.RemainQuantity,
				Rating:              template.Rating,
				NumberRating:        int32(template.NumberRating),
				SellerId:            memCacheProduct.SellerID,
				SellerName:          memCacheProduct.SellerName,
				SellerLogo:          memCacheProduct.SellerLogo.String,
				SellerAddress:       memCacheProduct.SellerAddress.String,
				CategoryId:          memCacheProduct.CategoryID,
				CategoryName:        memCacheProduct.CategoryName,
				UomId:               memCacheProduct.UomID,
				UomName:             memCacheProduct.UomName,
			},
		}, nil
	}

	// Get from redis

	// Get from database
	products, err := s.storeDb.GetProductDetails(ctx, []int64{req.Id})
	if err != nil {
		return nil, err
	}
	if len(products) == util.ZeroLength {
		return nil, fmt.Errorf("not found product with id = %v", req.Id)
	}
	data := &api.ProductDetail{}
	data.FromEntity(products[0])

	return &api.GetDetailProductResponse{
		Code:    0,
		Message: "OK",
		Data:    data,
	}, nil
}

// getOrInsertMem
//  Get Template from memCache
//  If mem not exists, get from database
//  If database exists, return
func getOrInsertMem(s *Service, ctx context.Context, memCacheProduct cache.Product) (cache.ProductTemplate, error) {
	template, ok := s.memCache.GetProductTemplate(memCacheProduct.TemplateID)
	if !ok {
		listTemplate, err := s.storeDb.GetProductTemplates(ctx, []int64{memCacheProduct.TemplateID})
		if err != nil {
			return cache.ProductTemplate{}, err
		}
		if len(listTemplate) == util.ZeroLength {
			return cache.ProductTemplate{}, status.Errorf(codes.NotFound, "Not found product(id = %v) with template id = %v", memCacheProduct.TemplateID, memCacheProduct.TemplateID)
		}
		template.FromDb(listTemplate[0])
	}
	return template, nil
}

func (s *Service) GetListProduct(ctx context.Context, req *api.GetListProductRequest) (res *api.GetListProductResponse, err error) {
	return nil, nil
}
