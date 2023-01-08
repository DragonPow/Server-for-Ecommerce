package service

import (
	"context"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
)

func (s *Service) GetDetailProduct(ctx context.Context, req *api.GetDetailProductRequest) (res *api.GetDetailProductResponse, err error) {
	// Get from memory

	// Get from redis

	// Get from database
	products, err := s.storeDb.GetProductDetails(ctx, []int64{req.Id})
	if err != nil {
		return nil, err
	}
	if len(products) == util.ZeroLength {
		return nil, fmt.Errorf("not found product with id = %v", req.Id)
	}
	product := products[0]
	return &api.GetDetailProductResponse{
		Code:    0,
		Message: "OK",
		Data: &api.ProductDetail{
			Id:                  product.ID,
			Name:                product.Name,
			OriginPrice:         product.OriginPrice,
			SalePrice:           product.SalePrice,
			Variants:            string(product.Variants.RawMessage),
			CreatedBy:           "",
			CreatedDate:         util.ParseTimeToString(product.CreateDate),
			UpdatedBy:           "",
			UpdatedDate:         util.ParseTimeToString(product.WriteDate),
			TemplateId:          product.TemplateID.Int64,
			TemplateName:        product.TemplateName,
			TemplateDescription: product.TemplateDescription.String,
			SoldQuantity:        product.SoldQuantity,
			RemainQuantity:      product.RemainQuantity,
			Rating:              product.Rating,
			NumberRating:        int32(product.NumberRating),
			SellerId:            product.SellerID,
			SellerName:          product.SellerName,
			SellerLogo:          product.SellerLogo.String,
			SellerAddress:       product.SellerAddress.String,
			CategoryId:          product.CategoryID,
			CategoryName:        product.CategoryName,
			UomId:               product.UomID,
			UomName:             product.UomName,
		},
	}, nil
}
func (s *Service) GetListProduct(ctx context.Context, req *api.GetListProductRequest) (res *api.GetListProductResponse, err error) {
	return nil, nil
}
