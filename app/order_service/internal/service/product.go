package service

import (
	"Server-for-Ecommerce/app/order_service/api"
	"Server-for-Ecommerce/app/order_service/internal/cache"
	"Server-for-Ecommerce/app/order_service/util"
	"context"
	"errors"
	"strconv"
)

func (s *Service) GetOrderDetail(ctx context.Context, request *api.GetOrderDetailRequest) (*api.GetOrderDetailResponse, error) {
	rsCache, ok := s.cache.Get(ctx, strconv.FormatInt(request.ProductId, util.Base10Int))
	if !ok {
		return nil, errors.New("Not found")
	}
	var res *api.GetOrderDetailResponse

	//if rsCache != nil {
	//	response, err := cache.UnMarshal[*api.GetOrderDetailResponse](rsCache)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return response, nil
	//}

	// Call db
	return nil, nil
}
