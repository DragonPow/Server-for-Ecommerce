package service

import (
	"context"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/database/store"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetDetailProduct(ctx context.Context, req *api.GetDetailProductRequest) (res *api.GetDetailProductResponse, err error) {
	logger := s.log.WithName("GetDetailProduct").WithValues("request", req)
	// Get from memory
	memCacheProduct, ok := s.memCache.GetProduct(req.Id)
	if ok {
		// Get Template
		template, err := getProductTemplateOrInsertCache(s, ctx, memCacheProduct.TemplateID)
		if err != nil {
			logger.Error(err, "getProductTemplateOrInsertCache")
			return nil, err
		}

		// Get Template
		category, err := getCategoryOrInsertCache(s, ctx, memCacheProduct.CategoryID)
		if err != nil {
			logger.Error(err, "getCategoryOrInsertCache")
			return nil, err
		}

		// Get Template
		uom, err := getUomOrInsertCache(s, ctx, memCacheProduct.UomID)
		if err != nil {
			logger.Error(err, "getUomOrInsertCache")
			return nil, err
		}

		// Get Template
		seller, err := getSellerOrInsertCache(s, ctx, memCacheProduct.SellerID)
		if err != nil {
			logger.Error(err, "getSellerOrInsertCache")
			return nil, err
		}

		// Get Template
		createBy, err := getUserOrInsertCache(s, ctx, memCacheProduct.CreateUid)
		if err != nil {
			logger.Error(err, "getUserOrInsertCache")
			return nil, err
		}

		// Get Template
		writeBy, err := getUserOrInsertCache(s, ctx, memCacheProduct.WriteUid)
		if err != nil {
			logger.Error(err, "getUserOrInsertCache")
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
				CreatedBy:           createBy.Name,
				CreatedDate:         util.ParseTimeToString(memCacheProduct.CreateDate),
				UpdatedBy:           writeBy.Name,
				UpdatedDate:         util.ParseTimeToString(memCacheProduct.WriteDate),
				TemplateId:          memCacheProduct.TemplateID,
				TemplateName:        template.Name,
				TemplateDescription: template.Description,
				SoldQuantity:        template.SoldQuantity,
				RemainQuantity:      template.RemainQuantity,
				Rating:              template.Rating,
				NumberRating:        int32(template.NumberRating),
				SellerId:            seller.ID,
				SellerName:          seller.Name,
				SellerLogo:          seller.Logo,
				SellerAddress:       seller.Address,
				CategoryId:          category.ID,
				CategoryName:        category.Name,
				UomId:               uom.ID,
				UomName:             uom.Name,
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

// getProductTemplateOrInsertCache
//  Get Template from memCache
//  If mem not exists, get from database
//  If database exists, return
func getProductTemplateOrInsertCache(s *Service, ctx context.Context, id int64) (template cache.ProductTemplate, err error) {
	var ok bool
	// Get from mem cache
	template, ok = s.memCache.GetProductTemplate(id)
	if !ok {
		// Get from redis
		template, ok = s.localCache.GetProductTemplate(id)
		if !ok {
			// Get from db
			listTemplate, err := s.storeDb.GetProductTemplates(ctx, []int64{id})
			if err != nil {
				return cache.ProductTemplate{}, err
			}
			if len(listTemplate) == util.ZeroLength {
				return cache.ProductTemplate{}, status.Errorf(codes.NotFound, "Not found template with id = %v", id)
			}
			template.FromDb(listTemplate[0])

			// Set to redis
			err = s.localCache.SetMultiple(map[int64]cache.ModelValue{id: template})
			if err != nil {
				s.log.Error(err, "Fail set multiple to local cache", "id", id, "template", template)
			}
		}

		// Check and set to mem cache
		_, err = s.memCache.CheckAndSet(map[int64]cache.ModelValue{id: template})
		if err != nil {
			s.log.Error(err, "Fail set multiple to mem cache", "id", id, "template", template)
		}
	}
	return template, nil
}

func getCategoryOrInsertCache(s *Service, ctx context.Context, id int64) (category cache.Category, err error) {
	var ok bool
	// Get from mem cache
	category, ok = s.memCache.GetCategory(id)
	if !ok {
		// Get from redis
		category, ok = s.localCache.GetCategory(id)
		if !ok {
			// Get from db
			listTemplate, err := s.storeDb.GetCategories(ctx, []int64{id})
			if err != nil {
				return cache.Category{}, err
			}
			if len(listTemplate) == util.ZeroLength {
				return cache.Category{}, status.Errorf(codes.NotFound, "Not found category with id = %v", id)
			}
			category.FromDb(listTemplate[0])

			// Set to redis
			err = s.localCache.SetMultiple(map[int64]cache.ModelValue{id: category})
			if err != nil {
				s.log.Error(err, "Fail set multiple to local cache", "id", id, "category", category)
			}
		}

		// Check and set to mem cache
		_, err = s.memCache.CheckAndSet(map[int64]cache.ModelValue{id: category})
		if err != nil {
			s.log.Error(err, "Fail set multiple to mem cache", "id", id, "category", category)
		}
	}
	return category, nil
}

func getUomOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.Uom, err error) {
	var (
		memUoms   map[int64]cache.Uom
		localUoms map[int64]cache.Uom
		dbUoms    map[int64]cache.Uom

		missMemIds   []int64
		missLocalIds []int64
	)
	// Get from mem cache
	memUoms, missMemIds = s.memCache.GetListUom(ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		localUoms, missLocalIds = s.localCache.GetListUom(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeUoms, err := s.storeDb.GetUoms(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeUoms) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found uoms with ids = %v", missLocalIds)
			}
			dbUoms = math.ToMap(storeUoms, func(uom store.Uom) (int64, cache.Uom) {
				var u cache.Uom
				u.FromDb(uom)
				return uom.ID, u
			})

			// Set to redis
			err = s.localCache.SetMultiple(math.ConvertMap(dbUoms, util.FuncConvertToCache[cache.Uom]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "dbUoms", dbUoms)
			}
		}

		// Check and set to mem cache
		newMemCache := make(map[int64]cache.Uom, len(dbUoms)+len(localUoms))
		for id, uom := range localUoms {
			newMemCache[id] = uom
		}
		for id, uom := range dbUoms {
			newMemCache[id] = uom
		}
		_, err = s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[cache.Uom]))
		if err != nil {
			s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "localUoms", localUoms, "dbUoms", dbUoms)
		}
	}

	result = make(map[int64]cache.Uom, len(ids))
	for id, uom := range memUoms {
		result[id] = uom
	}
	for id, uom := range localUoms {
		result[id] = uom
	}
	for id, uom := range dbUoms {
		result[id] = uom
	}
	return result, nil
}

func getOrInsertCache[T cache.ModelValue](s *Service, ctx context.Context, ids []int64) (result map[int64]T, err error) {
	var (
		mems   map[int64]T
		locals map[int64]T
		dbs    map[int64]T

		missMemIds   []int64
		missLocalIds []int64
	)
	// Get from mem cache
	mems, missMemIds = s.memCache.GetList(*new(T), ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		locals, missLocalIds = s.localCache.GetList(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeUoms, err := s.storeDb.GetUoms(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeUoms) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found uoms with ids = %v", missLocalIds)
			}
			dbs = math.ToMap(storeUoms, func(uom store.Uom) (int64, cache.Uom) {
				var u cache.Uom
				u.FromDb(uom)
				return uom.ID, u
			})

			// Set to redis
			err = s.localCache.SetMultiple(math.ConvertMap(dbUoms, util.FuncConvertToCache[cache.Uom]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "dbUoms", dbUoms)
			}
		}

		// Check and set to mem cache
		newMemCache := make(map[int64]cache.Uom, len(dbUoms)+len(localUoms))
		for id, uom := range localUoms {
			newMemCache[id] = uom
		}
		for id, uom := range dbUoms {
			newMemCache[id] = uom
		}
		_, err = s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[cache.Uom]))
		if err != nil {
			s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "localUoms", localUoms, "dbUoms", dbUoms)
		}
	}

	result = make(map[int64]cache.Uom, len(ids))
	for id, uom := range memUoms {
		result[id] = uom
	}
	for id, uom := range localUoms {
		result[id] = uom
	}
	for id, uom := range dbUoms {
		result[id] = uom
	}
	return result, nil
}

func (s *Service) GetListProduct(ctx context.Context, req *api.GetListProductRequest) (res *api.GetListProductResponse, err error) {
	return nil, nil
}
