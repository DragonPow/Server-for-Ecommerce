package service

import (
	"Server-for-Ecommerce/app_v2/product_service/api"
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"context"
	"errors"
)

func (s *Service) DeleteCache(ctx context.Context, req *api.DeleteCacheRequest) (res *api.DeleteCacheResponse, err error) {
	logger := s.log.WithName("DeleteCache").WithValues("req", req)
	switch req.Level {
	case "mem":
		err := s.memCache.Delete(cache.TypeCache(req.TypeModel), req.Ids)
		if err != nil {
			return nil, err
		}
		logger.Info("Delete mem success")
		return &api.DeleteCacheResponse{
			Code:    0,
			Message: "OK",
		}, nil
	case "redis":
		err := s.localCache.Delete(cache.TypeCache(req.TypeModel), req.Ids)
		if err != nil {
			return nil, err
		}
		logger.Info("Delete local success")
		return &api.DeleteCacheResponse{
			Code:    0,
			Message: "OK",
		}, nil
	case "page":
		err = s.localCache.DeletePage(req.Page, req.PageSize, req.Key)
		if err != nil {
			return nil, err
		}
		logger.Info("Delete page product success")
		return &api.DeleteCacheResponse{
			Code:    0,
			Message: "OK",
		}, nil
	default:
		return nil, errors.New("Level must be mem or  cache")
	}
}
