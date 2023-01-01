package service

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
)

func (s *Service) GetDetailProduct(ctx context.Context, req *api.GetDetailProductRequest) (res *api.GetDetailProductResponse, err error) {
	// Get from memory
	// Get from database
	return &api.GetDetailProductResponse{
		Code:    0,
		Message: "OK",
		Data: api.ProductDetail{
			Id:                  100,
			Name:                "Máy tính xách tay ACER Nitro 5",
			OriginPrice:         30000000,
			SalePrice:           45000000,
			Variants:            `{"màu":"đỏ/đen", "cân nặng":"15kg", "đơn vị giao hàng":"Giao hàng tiết kiệm","bảo hành":"12 tháng"}`,
			CreatedBy:           "Thạch Vũ Ngọc",
			CreatedDate:         "2011-10-05T14:48:00.000Z",
			UpdatedBy:           "Thạch Vũ Ngọc",
			UpdatedDate:         "2011-10-05T14:48:00.000Z",
			Decription:          "Máy tính chuyên dụng chơi game văn phòng hoặc sử dụng cho sinh viên làm bài tập, chơi game",
			TemplateId:          1,
			TemplateName:        "",
			TemplateDescription: "",
			SoldQuantity:        7000,
			RemainQuantity:      1024,
			Rating:              3.5,
			NumberRating:        4000,
			SellerId:            1,
			SellerName:          "Công ty trách nhiệm hữu hạn 1 thành viên",
			SellerLogo:          "",
			SellerAddress:       "Q1, TP.HCM",
			CategoryId:          1,
			CategoryName:        "Máy tính",
			UomId:               1,
			UomName:             "Cái",
		},
	}, nil
}
func (s *Service) GetListProduct(ctx context.Context, req *api.GetListProductRequest) (res *api.GetListProductResponse, err error) {
	return nil, nil
}
