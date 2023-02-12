package service

import (
	"Server-for-Ecommerce/app_v2/product_service/api"
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/app_v2/product_service/database/store"
	"Server-for-Ecommerce/app_v2/product_service/util"
	"Server-for-Ecommerce/library/math"
	"Server-for-Ecommerce/library/slice"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
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
		return s.computeGetDetailProduct(ctx, logger, memCacheProduct)
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
		return s.computeGetDetailProduct(ctx, logger, localCacheProduct)
	}

	type funcCallDb struct {
		wg  *sync.WaitGroup
		res *api.GetDetailProductResponse
		err error
	}

	keyLoad := fmt.Sprintf("product:%d", req.Id)
	s.lockCache.mu.Lock()
	v, ok := s.lockCache.list.Load(keyLoad)
	if ok {
		rs := v.(*funcCallDb)
		s.lockCache.mu.Unlock()
		logger.Info("Wait lockCache")
		rs.wg.Wait()
		return rs.res, rs.err
	} else {
		rs := &funcCallDb{
			wg:  &sync.WaitGroup{},
			res: nil,
			err: nil,
		}
		rs.wg.Add(1)
		s.lockCache.list.Store(keyLoad, rs)
		s.lockCache.mu.Unlock()
		// Get from database
		rs.res, rs.err = s.FetchFromDb(ctx, req, logger)
		rs.wg.Done()
		return rs.res, rs.err
	}

}

func (s *Service) FetchFromDb(ctx context.Context, req *api.GetDetailProductRequest, logger logr.Logger) (*api.GetDetailProductResponse, error) {
	products, err := s.storeDb.GetProductDetails(ctx, []int64{req.Id})
	if err != nil {
		logger.Error(err, "GetProductDetails fail")
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

func (s *Service) computeGetDetailProduct(ctx context.Context, logger logr.Logger, cacheModel cache.Product) (*api.GetDetailProductResponse, error) {
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
		b, _ := json.Marshal(cacheModel.Variants)
		return &api.GetDetailProductResponse{
			Code:    0,
			Message: "OK",
			Data: &api.ProductDetail{
				Id:                  cacheModel.ID,
				Name:                cacheModel.Name,
				OriginPrice:         cacheModel.OriginPrice,
				SalePrice:           cacheModel.SalePrice,
				Variants:            string(b),
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
				Image:               cacheModel.Image,
			},
		}, nil
	}
}

func getProductOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.Product, err error) {
	logger := s.log.WithName("getProductOrInsertCache")
	type typeCache = cache.Product
	type typeDb = store.Product
	var (
		modelName        = "product"
		funcMemGetList   = s.memCache.GetListProduct
		funcLocalGetList = s.localCache.GetListProduct
		funcDbGetList    = s.storeDb.GetProductAndRelations

		mem   map[int64]typeCache
		local map[int64]typeCache
		db    map[int64]typeCache

		missMemIds   []int64
		missLocalIds []int64
	)

	// Get from mem cache
	mem, missMemIds = funcMemGetList(ids)
	logger.Info("Get from mem cache", "mem", maps.Keys(mem))
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		logger.Info("Get from local cache", "local", maps.Keys(local))
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) < len(missLocalIds) {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds, "store", storeModel)
			}
			logger.Info("Get from db", "store", missLocalIds)
			db = math.ToMap(storeModel, func(model store.GetProductAndRelationsRow) (int64, typeCache) {
				var u typeCache
				u.FromDbV2(model)
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

func getProductTemplateOrInsertCache(s *Service, ctx context.Context, ids []int64) (result map[int64]cache.ProductTemplate, err error) {
	logger := s.log.WithName("getProductTemplateOrInsertCache")
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
	logger.Info("Get from mem cache", "mem", maps.Keys(mem))
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		logger.Info("Get from local cache", "local", maps.Keys(local))
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) < len(missLocalIds) {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds, "store", storeModel)
			}
			logger.Info("Get from db", "store", missLocalIds)
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
	logger := s.log.WithName("getUomOrInsertCache")
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
	logger.Info("Get from mem cache", "mem", maps.Keys(mem))
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		logger.Info("Get from local cache", "local", maps.Keys(local))
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) < len(missLocalIds) {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds, "store", storeModel)
			}
			logger.Info("Get from db", "store", missLocalIds)
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
	logger := s.log.WithName("getCategoryOrInsertCache")
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
	logger.Info("Get from mem cache", "mem", maps.Keys(mem))
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		logger.Info("Get from local cache", "local", maps.Keys(local))
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) < len(missLocalIds) {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds, "store", storeModel)
			}
			logger.Info("Get from db", "store", missLocalIds)
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
	logger := s.log.WithName("getSellerOrInsertCache")
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
	logger.Info("Get from mem cache", "mem", maps.Keys(mem))
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		logger.Info("Get from local cache", "local", maps.Keys(local))
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) < len(missLocalIds) {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds, "store", storeModel)
			}
			logger.Info("Get from db", "store", missLocalIds)
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
	logger := s.log.WithName("getUserOrInsertCache")
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
	logger.Info("Get from mem cache", "mem", maps.Keys(mem))
	if len(missMemIds) > util.ZeroLength {
		// Get from redis
		local, missLocalIds = funcLocalGetList(missMemIds)
		logger.Info("Get from local cache", "local", maps.Keys(local))
		if len(missLocalIds) > util.ZeroLength {
			// Get from db
			storeModel, err := funcDbGetList(ctx, missLocalIds)
			if err != nil {
				return nil, err
			}
			if len(storeModel) < len(missLocalIds) {
				return nil, status.Errorf(codes.NotFound, "Not found %s with ids = %v", modelName, missLocalIds, "store", storeModel)
			}
			logger.Info("Get from db", "store", missLocalIds)
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
	logger := s.log.WithName("GetListProduct").WithValues("request", req)
	isCacheRedisPage := req.Page < s.cfg.RedisConfig.NumberCachePage
	if isCacheRedisPage {
		res = &api.GetListProductResponse{}
		pageCache, ok := s.localCache.GetPageProduct(req.Page, req.PageSize, req.Key)
		if ok {
			err := json.Unmarshal([]byte(pageCache), res)
			if err != nil {
				logger.Error(err, "Unmarshall pageCache fail")
				return nil, err
			}
			logger.Info("Get page product from cache")
			return res, nil
		}
	}

	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	rows, err := s.storeDb.GetProductsByKeyword(ctx, store.GetProductsByKeywordParams{
		Keyword: fmt.Sprintf("%%%v%%", req.Key),
		Offset:  offset,
		Limit:   limit,
	})
	if err != nil {
		logger.Error(err, "GetProductsByKeyword fail")
		return nil, err
	}
	res = &api.GetListProductResponse{
		Code:    0,
		Message: "OK",
		Data: &api.GetListProductResponse_Data{
			TotalItems: 0,
			Page:       int32(req.Page),
			PageSize:   int32(req.PageSize),
			Items:      nil,
		},
	}

	if isCacheRedisPage {
		go func(logger logr.Logger) {
			data, err := json.Marshal(res)
			if err != nil {
				logger.Error(err, "Marshal res fail")
				return
			}
			err = s.localCache.SetPageProduct(
				req.Page, req.PageSize, req.Key,
				string(data),
				time.Duration(s.cfg.RedisConfig.ExpireCachePageInSecond)*time.Second,
			)
			if err != nil {
				logger.Error(err, "Set page product fail")
				return
			}
			logger.Info("Set page product success")
		}(logger)

	}

	if len(rows) == util.ZeroLength {
		logger.Info("Not found items")
		return res, nil
	}

	totalItems := rows[0].Total
	items := math.Convert(rows, func(row store.GetProductsByKeywordRow) *api.ProductOverview {
		return &api.ProductOverview{
			Id:          row.ID,
			Name:        row.Name,
			OriginPrice: row.OriginPrice,
			SalePrice:   row.SalePrice,
			Image:       row.Image,
		}
	})
	logger.Info("Get list product success", "ids", slice.Map(rows, func(row store.GetProductsByKeywordRow) int64 { return row.ID }))
	return &api.GetListProductResponse{
		Code:    0,
		Message: "OK",
		Data: &api.GetListProductResponse_Data{
			TotalItems: totalItems,
			Page:       int32(req.Page),
			PageSize:   int32(req.PageSize),
			Items:      items,
		},
	}, nil
}
