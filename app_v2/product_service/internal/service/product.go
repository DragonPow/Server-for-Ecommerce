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
		templates, err := getProductTemplateOrInsertCache(s, ctx, []int64{memCacheProduct.TemplateID})
		if err != nil {
			logger.Error(err, "getProductTemplateOrInsertCache")
			return nil, err
		}
		template := templates[memCacheProduct.TemplateID]

		// Get Template
		categories, err := getCategoryOrInsertCache(s, ctx, []int64{memCacheProduct.CategoryID})
		if err != nil {
			logger.Error(err, "getCategoryOrInsertCache")
			return nil, err
		}
		category := categories[memCacheProduct.CategoryID]

		// Get Template
		uoms, err := getUomOrInsertCache(s, ctx, []int64{memCacheProduct.UomID})
		if err != nil {
			logger.Error(err, "getUomOrInsertCache")
			return nil, err
		}
		uom := uoms[memCacheProduct.UomID]

		// Get Template
		sellers, err := getSellerOrInsertCache(s, ctx, []int64{memCacheProduct.SellerID})
		if err != nil {
			logger.Error(err, "getSellerOrInsertCache")
			return nil, err
		}
		seller := sellers[memCacheProduct.SellerID]

		// Get Template
		createBys, err := getUserOrInsertCache(s, ctx, []int64{memCacheProduct.CreateUid})
		if err != nil {
			logger.Error(err, "getUserOrInsertCache")
			return nil, err
		}
		createBy := createBys[memCacheProduct.CreateUid]

		// Get Template
		writeBys, err := getUserOrInsertCache(s, ctx, []int64{memCacheProduct.WriteUid})
		if err != nil {
			logger.Error(err, "getUserOrInsertCache")
			return nil, err
		}
		writeBy := writeBys[memCacheProduct.WriteUid]

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

func getProductTemplateOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.ProductTemplate, err error) {
	type typeCache = cache.ProductTemplate
	type typeDb = store.ProductTemplate
	var (
		modelName        = "uom"
		funcMemGetList   = s.memCache.GetListProductTemplate
		funcLocalGetList = s.localCache.GetListProductTemplate
		funcDbGetList    = s.storeDb.GetProductTemplates

		mem   map[int64]typeCache
		local map[int64]typeCache
		db    map[int64]typeCache

		missMemIds   []int64
		missLocalIds []int64
	)

	// Get from mem cache
	mem, missMemIds = funcMemGetList(ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds)
			}
			db = math.ToMap(storeModel, func(model typeDb) (int64, typeCache) {
				var u typeCache
				u.FromDb(model)
				return model.ID, u
			})

			// Set to redis
			go func() {
				err := s.localCache.SetMultiple(math.ConvertMap(db, util.FuncConvertToCache[typeCache]))
				if err != nil {
					s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "db", db)
				}
			}()
		}

		// Check and set to mem cache
		go func() {
			newMemCache := math.AppendMap(db, local)
			_, err := s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[typeCache]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "local", local, "db", db)
			}
		}()
	}

	return math.AppendMap(mem, local, db), nil
}

func getUomOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.Uom, err error) {
	type typeCache = cache.Uom
	type typeDb = store.Uom
	var (
		modelName        = "uom"
		funcMemGetList   = s.memCache.GetListUom
		funcLocalGetList = s.localCache.GetListUom
		funcDbGetList    = s.storeDb.GetUoms

		mem   map[int64]typeCache
		local map[int64]typeCache
		db    map[int64]typeCache

		missMemIds   []int64
		missLocalIds []int64
	)

	// Get from mem cache
	mem, missMemIds = funcMemGetList(ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds)
			}
			db = math.ToMap(storeModel, func(model typeDb) (int64, typeCache) {
				var u typeCache
				u.FromDb(model)
				return model.ID, u
			})

			// Set to redis
			go func() {
				err := s.localCache.SetMultiple(math.ConvertMap(db, util.FuncConvertToCache[typeCache]))
				if err != nil {
					s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "db", db)
				}
			}()
		}

		// Check and set to mem cache
		go func() {
			newMemCache := math.AppendMap(db, local)
			_, err := s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[typeCache]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "local", local, "db", db)
			}
		}()
	}

	return math.AppendMap(mem, local, db), nil
}

func getCategoryOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.Category, err error) {
	type typeCache = cache.Category
	type typeDb = store.Category
	var (
		modelName        = "uom"
		funcMemGetList   = s.memCache.GetListCategory
		funcLocalGetList = s.localCache.GetListCategory
		funcDbGetList    = s.storeDb.GetCategories

		mem   map[int64]typeCache
		local map[int64]typeCache
		db    map[int64]typeCache

		missMemIds   []int64
		missLocalIds []int64
	)

	// Get from mem cache
	mem, missMemIds = funcMemGetList(ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds)
			}
			db = math.ToMap(storeModel, func(model typeDb) (int64, typeCache) {
				var u typeCache
				u.FromDb(model)
				return model.ID, u
			})

			// Set to redis
			go func() {
				err := s.localCache.SetMultiple(math.ConvertMap(db, util.FuncConvertToCache[typeCache]))
				if err != nil {
					s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "db", db)
				}
			}()
		}

		// Check and set to mem cache
		go func() {
			newMemCache := math.AppendMap(db, local)
			_, err := s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[typeCache]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "local", local, "db", db)
			}
		}()
	}

	return math.AppendMap(mem, local, db), nil
}

func getSellerOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.Seller, err error) {
	type typeCache = cache.Seller
	type typeDb = store.Seller
	var (
		modelName        = "uom"
		funcMemGetList   = s.memCache.GetListSeller
		funcLocalGetList = s.localCache.GetListSeller
		funcDbGetList    = s.storeDb.GetSellers

		mem   map[int64]typeCache
		local map[int64]typeCache
		db    map[int64]typeCache

		missMemIds   []int64
		missLocalIds []int64
	)

	// Get from mem cache
	mem, missMemIds = funcMemGetList(ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds)
			}
			db = math.ToMap(storeModel, func(model typeDb) (int64, typeCache) {
				var u typeCache
				u.FromDb(model)
				return model.ID, u
			})

			// Set to redis
			go func() {
				err := s.localCache.SetMultiple(math.ConvertMap(db, util.FuncConvertToCache[typeCache]))
				if err != nil {
					s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "db", db)
				}
			}()
		}

		// Check and set to mem cache
		go func() {
			newMemCache := math.AppendMap(db, local)
			_, err := s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[typeCache]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "local", local, "db", db)
			}
		}()
	}

	return math.AppendMap(mem, local, db), nil
}

func getUserOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.User, err error) {
	type typeCache = cache.User
	type typeDb = store.User
	var (
		modelName        = "uom"
		funcMemGetList   = s.memCache.GetListUser
		funcLocalGetList = s.localCache.GetListUser
		funcDbGetList    = s.storeDb.GetUsers

		mem   map[int64]typeCache
		local map[int64]typeCache
		db    map[int64]typeCache

		missMemIds   []int64
		missLocalIds []int64
	)

	// Get from mem cache
	mem, missMemIds = funcMemGetList(ids)
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) == util.ZeroLength {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds)
			}
			db = math.ToMap(storeModel, func(model typeDb) (int64, typeCache) {
				var u typeCache
				u.FromDb(model)
				return model.ID, u
			})

			// Set to redis
			go func() {
				err := s.localCache.SetMultiple(math.ConvertMap(db, util.FuncConvertToCache[typeCache]))
				if err != nil {
					s.log.Error(err, "Fail set multiple to local cache", "ids", ids, "db", db)
				}
			}()
		}

		// Check and set to mem cache
		go func() {
			newMemCache := math.AppendMap(db, local)
			_, err := s.memCache.CheckAndSet(math.ConvertMap(newMemCache, util.FuncConvertToCache[typeCache]))
			if err != nil {
				s.log.Error(err, "Fail set multiple to mem cache", "ids", ids, "local", local, "db", db)
			}
		}()
	}

	return math.AppendMap(mem, local, db), nil
}

func (s *Service) GetListProduct(ctx context.Context, req *api.GetListProductRequest) (res *api.GetListProductResponse, err error) {
	return nil, nil
}
