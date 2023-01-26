package service

import (
	"context"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/cache"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/database/store"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

func (s *Service) GetDetailProduct(ctx context.Context, req *api.GetDetailProductRequest) (res *api.GetDetailProductResponse, err error) {
	logger := s.log.WithName("GetDetailProduct").WithValues("request", req)
	var (
		memCacheProduct   cache.Product
		localCacheProduct cache.Product
		hasCache          bool
	)

	// Get from memory cache
	memCacheProduct, hasCache = s.memCache.GetProduct(req.Id)
	if hasCache {
		logger.Info("Get data from mem cache")
		return s.computeFromCache(ctx, logger, memCacheProduct)
	}

	// Get from redis
	localCacheProduct, hasCache = s.localCache.GetProduct(req.Id)
	if hasCache {
		defer func() {
			go func() {
				// Check and set to mem cache
				ok, err := s.memCache.CheckAndSet(map[int64]cache.ModelValue{req.Id: localCacheProduct})
				if err != nil {
					s.log.Error(err, "Fail set product to mem cache", "id", req.Id, "local", localCacheProduct)
				}
				if ok {
					s.log.Info("Set multiple mem cache success", "productId", req.Id)
				}
			}()
		}()
		logger.Info("Get data from local cache")
		return s.computeFromCache(ctx, logger, localCacheProduct)
	}

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
	// Update cache for local and mem
	go s.SetCacheAPIGetDetailProduct(data)

	logger.Info("Get data from db")
	return &api.GetDetailProductResponse{
		Code:    0,
		Message: "OK",
		Data:    data,
	}, nil
}

func (s *Service) SetCacheAPIGetDetailProduct(data *api.ProductDetail) {
	ctx := context.Background()
	p, err := s.storeDb.GetProducts(ctx, []int64{data.Id})
	if err != nil {
		s.log.Error(err, "GetProducts", "id", data.Id)
		return
	}
	if len(p) == util.ZeroLength {
		// Ignore
		s.log.Info("Not found product to set to local, mem", "id", data.Id)
		return
	}
	productCache := cache.Product{}
	productCache.FromDb(p[0], data.CategoryId, data.UomId, data.SellerId)
	// Insert redis cache
	go func() {
		err := s.localCache.SetMultiple(map[int64]cache.ModelValue{data.Id: productCache})
		if err != nil {
			s.log.Error(err, "Set multiple local cache fail")
			return
		}
		s.log.Info("Set multiple local cache success", "productId", data.Id)
	}()
	// Insert mem cache
	go func() {
		ok, err := s.memCache.CheckAndSet(map[int64]cache.ModelValue{data.Id: productCache})
		if err != nil {
			s.log.Error(err, "Set multiple mem cache fail")
			return
		}
		if ok {
			s.log.Info("Set multiple mem cache success", "productId", data.Id)
		}
	}()
}

func (s *Service) computeFromCache(ctx context.Context, logger logr.Logger, cacheModel cache.Product) (*api.GetDetailProductResponse, error) {
	wg := &sync.WaitGroup{}
	errChan := make(chan error)
	doneChan := make(chan struct{})
	wg.Add(5)
	var (
		template cache.ProductTemplate
		uom      cache.Uom
		category cache.Category
		seller   cache.Seller
		createBy cache.User
		writeBy  cache.User
	)
	//var (
	//	templateChan chan cache.ProductTemplate
	//	uomChan      chan cache.Uom
	//	categoryChan chan cache.Category
	//	sellerChan   chan cache.Seller
	//	createByChan chan cache.User
	//	writeByChan  chan cache.User
	//)

	// Get Template
	go func() {
		defer wg.Done()
		templates, err := getProductTemplateOrInsertCache(s, ctx, []int64{cacheModel.TemplateID})
		if err != nil {
			logger.Error(err, "getProductTemplateOrInsertCache")
			errChan <- err
			return
		}
		template = templates[cacheModel.TemplateID]
	}()

	// Get Category
	go func() {
		defer wg.Done()
		categories, err := getCategoryOrInsertCache(s, ctx, []int64{cacheModel.CategoryID})
		if err != nil {
			logger.Error(err, "getCategoryOrInsertCache")
			errChan <- err
		}
		category = categories[cacheModel.CategoryID]
	}()

	// Get Uom
	go func() {
		defer wg.Done()
		uoms, err := getUomOrInsertCache(s, ctx, []int64{cacheModel.UomID})
		if err != nil {
			logger.Error(err, "getUomOrInsertCache")
			errChan <- err
		}
		uom = uoms[cacheModel.UomID]
	}()

	// Get Seller
	go func() {
		defer wg.Done()
		sellers, err := getSellerOrInsertCache(s, ctx, []int64{cacheModel.SellerID})
		if err != nil {
			logger.Error(err, "getSellerOrInsertCache")
			errChan <- err
		}
		seller = sellers[cacheModel.SellerID]
	}()

	// Get CreateBy
	go func() {
		defer wg.Done()
		users, err := getUserOrInsertCache(s, ctx, math.Uniq([]int64{cacheModel.CreateUid, cacheModel.WriteUid}))
		if err != nil {
			logger.Error(err, "getUserOrInsertCache")
			errChan <- err
		}
		createBy = users[cacheModel.CreateUid]
		writeBy = users[cacheModel.WriteUid]
	}()

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-doneChan:

		return &api.GetDetailProductResponse{
			Code:    0,
			Message: "OK",
			Data: &api.ProductDetail{
				Id:                  cacheModel.ID,
				Name:                cacheModel.Name,
				OriginPrice:         cacheModel.OriginPrice,
				SalePrice:           cacheModel.SalePrice,
				Variants:            cacheModel.Variants,
				CreatedBy:           createBy.Name,
				CreatedDate:         util.ParseTimeToString(cacheModel.CreateDate),
				UpdatedBy:           writeBy.Name,
				UpdatedDate:         util.ParseTimeToString(cacheModel.WriteDate),
				TemplateId:          cacheModel.TemplateID,
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
}

func getProductTemplateOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.ProductTemplate, err error) {
	type typeCache = cache.ProductTemplate
	type typeDb = store.ProductTemplate
	var (
		modelName        = "productTemplate"
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
		modelName        = "category"
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
		modelName        = "seller"
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
		modelName        = "user"
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
