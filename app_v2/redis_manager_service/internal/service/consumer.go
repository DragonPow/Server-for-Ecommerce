package service

import (
	producerDb "Server-for-Ecommerce/app_v2/db_manager_service/util"
	"Server-for-Ecommerce/app_v2/product_service/cache"
	storeProduct "Server-for-Ecommerce/app_v2/product_service/database/store"
	"Server-for-Ecommerce/app_v2/redis_manager_service/internal/database/store"
	"Server-for-Ecommerce/app_v2/redis_manager_service/util"
	producer "Server-for-Ecommerce/library/kafka/pub"
	"Server-for-Ecommerce/library/ring"
	"Server-for-Ecommerce/library/slice"
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"golang.org/x/exp/maps"
	"io"
	"time"
)

func (s *Service) Consume() error {
	kafkaConfig := s.cfg.KafkaConfig
	errChan := make(chan error)

	// Wait ring signal
	if s.cfg.EnableRing {
		go func() {
			for {
				if s.redis.Ring.Length() == util.ZeroLength {
					//s.log.Info("Ring empty, wait write")
					<-s.redis.Ring.SignalWrite
					//s.log.Info("Write success, continue")
				}
				timeout := time.After(time.Duration(s.redis.TimeoutRingWriterInMillisecond) * time.Millisecond)
				err := s.WaitRing(timeout)
				if err != nil {
					//s.log.Error(err, "WaitRing fail")
					continue
				}
				//s.log.Info("Wait ring success")
			}
		}()
	}

	var funcConsumeUpdate func(ctx context.Context, message kafka.Message) error
	if s.cfg.EnableRing {
		funcConsumeUpdate = s.AddUpdateMessageToRing
	} else {
		funcConsumeUpdate = s.computeMessagesUpdateV3
	}

	// Consumer update database
	updateConsumer := kafkaConfig.UpdateDbConsumer
	go func() {
		// create a new reader to the topic "update-db"
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:                updateConsumer.Connections,
			Topic:                  updateConsumer.Topic,
			GroupID:                updateConsumer.Group,
			MinBytes:               10e3, // 10KB
			MaxBytes:               10e6, // 10MB
			StartOffset:            kafka.LastOffset,
			MaxWait:                1 * time.Second,
			PartitionWatchInterval: 1 * time.Second,
		})
		err := s.ProcessConsume(r, funcConsumeUpdate)
		if err != nil {
			errChan <- err
		}
	}()

	//go func() {
	//	// create a new reader to the topic "insert-db"
	//	r := kafka.NewReader(kafka.ReaderConfig{
	//		Brokers:  updateConsumer.Connections,
	//		Topic:    updateConsumer.Topic,
	//		GroupID:  updateConsumer.Group,
	//		MinBytes: 10e3, // 10KB
	//		MaxBytes: 10e6, // 10MB
	//	})
	//	err := s.ProcessConsume(r, s.UpdateRedis)
	//	if err != nil {
	//		errChan <- err
	//	}
	//}()

	return <-errChan
}

func (s *Service) ProcessConsume(r *kafka.Reader, process func(ctx context.Context, message kafka.Message) error) error {
	defer r.Close()
	// consume messages
	for {
		ctx := context.Background()
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				//s.log.Info("Message EOF")
				continue
			}
			s.log.Error(err, "Read message fail")
			return err
		}

		logger := s.log.WithValues("Offset", m.Offset, "Key", string(m.Key), "Value", string(m.Value))
		logger.Info("Receive message")
		go func() {
			err := process(ctx, m)
			if err != nil {
				logger.Error(err, "Process consumer fail")
			}
		}()
	}
}

func (s *Service) WaitRing(timeout <-chan time.Time) error {
	for {
		select {
		case <-s.redis.Ring.SignalFull:
			s.log.Info("SignalFull")
			return s.updateRedisFromRing()
		case <-timeout:
			s.log.Info("SignalTimeout")
			return s.updateRedisFromRing()
		case <-s.redis.Ring.SignalWrite:
		}
	}
}

func (s *Service) updateRedisFromRing() error {
	logger := s.log.WithName("updateRedisFromRing")
	if s.redis.Ring.IsEmpty() {
		return nil
	}
	length := s.redis.Ring.Length()
	buf := make([]kafka.Message, length)
	_, err := s.redis.Ring.Read(buf)
	if errors.Is(ring.ErrIsEmpty, err) {
		logger.Info("Ring is empty")
		return nil
	}
	if err != nil {
		logger.Error(err, "Read from ring fail")
		s.redis.Ring.Free()
		return err
	}
	// Compute all message in ring
	err = s.computeMessagesUpdateV2(buf)
	if err != nil {
		logger.Error(err, "computeMessagesUpdate fail")
		return err
	}
	logger.Info("Success")
	return nil
}

func (s *Service) AddUpdateMessageToRing(ctx context.Context, message kafka.Message) error {
	logger := s.log.WithName("AddUpdateMessageToRing").WithValues("message", message)
	err := s.redis.Ring.WriteOne(message)
	if err != nil {
		logger.Error(err, "Write fail")
		return err
	}
	logger.Info("Success")
	return nil
}

//func (s *Service) computeMessagesUpdate(buf []kafka.Message) error {
//	mapBuf := make(map[int64]*producerDb.UpdateDatabaseEventValue, len(buf))
//	for _, message := range buf {
//		// Marshal from message
//		var payload producerDb.UpdateDatabaseEventValue
//		err := json.Unmarshal(message.Value, &payload)
//		if err != nil {
//			return err
//		}
//
//		// Check if mapBuf exists id, append to variants
//		// Else add to mapBuf
//		m, ok := mapBuf[payload.Id]
//		if !ok {
//			mapBuf[payload.Id] = &payload
//			continue
//		}
//
//		// Compute variants
//		variants := make(map[string]any)
//		mVariants := make(map[string]any)
//		var isAppend bool
//		err = json.Unmarshal(payload.Variants, &variants)
//		if err != nil {
//			return err
//		}
//		err = json.Unmarshal(m.Variants, &mVariants)
//		if err != nil {
//			return err
//		}
//
//		// Merge prev_variants and payload Variants
//		// Only merge when old message
//		for k, v := range variants {
//			isAppend = true
//			_, ok := mVariants[k]
//			if ok {
//				// If exists in mapBuf and version is lower, change in mapBuf
//				if m.GetVersion() < payload.GetVersion() {
//					mVariants[k] = v
//				}
//				continue
//			}
//			mVariants[k] = v
//		}
//		// If merge, update in mapBuf variants
//		if isAppend {
//			b, err := json.Marshal(mVariants)
//			if err != nil {
//				return err
//			}
//			m.Variants = b
//		}
//		continue
//	}
//
//	for _, payload := range mapBuf {
//		go func(payload producerDb.UpdateDatabaseEventValue) {
//			ctx := context.Background()
//			_ = s.UpdateRedis(ctx, payload) // ignore if error
//		}(*payload)
//	}
//	return nil
//}

func (s *Service) computeMessagesUpdateV2(buf []kafka.Message) error {
	mapBuf := make(map[int64]struct{}, len(buf))
	for _, message := range buf {
		// Marshal from message
		var payload producerDb.UpdateDatabaseEventValue
		err := json.Unmarshal(message.Value, &payload)
		if err != nil {
			return err
		}

		// Check if mapBuf exists id, append to variants
		// Else add to mapBuf
		mapBuf[payload.Id] = struct{}{}
	}

	return s.UpdateRedisV2(context.Background(), maps.Keys(mapBuf))
}

func (s *Service) computeMessagesUpdateV3(ctx context.Context, message kafka.Message) error {
	// Marshal from message
	var payload producerDb.UpdateDatabaseEventValue
	err := json.Unmarshal(message.Value, &payload)
	if err != nil {
		return err
	}
	// Update
	return s.UpdateRedisV2(context.Background(), []int64{payload.Id})
}

//func (s *Service) UpdateRedis(ctx context.Context, payload producerDb.UpdateDatabaseEventValue) error {
//	logger := s.log.WithName("UpdateMemoryCache").WithValues("payload", payload)
//	logger.Info("Start process")
//
//	// Way 1: Call db to enrich data
//	//products, err := s.storeDb.GetProducts(ctx, []int64{payload.Id})
//	//if err != nil {
//	//	logger.Error(err, "Call db get products fail")
//	//}
//	//if len(products) == util.ZeroLength {
//	//	err = fmt.Errorf("not found product with id %v", payload.Id)
//	//	logger.Error(err, "Product not exists")
//	//	return err
//	//}
//	//product := products[0]
//	//templates, err := s.storeDb.GetProductTemplates(ctx, []int64{product.TemplateID.Int64})
//	//if err != nil {
//	//	logger.Error(err, "Call db get product templates fail")
//	//}
//	//if len(templates) == util.ZeroLength {
//	//	err = fmt.Errorf("not found product template with id %v", payload.Id)
//	//	logger.Error(err, "Product not exists")
//	//	return err
//	//}
//	//template := templates[0]
//	//var cacheModel cache.Product
//	//cacheModel.FromDb(storeProduct.Product(product), template.CategoryID.Int64, template.UomID.Int64, template.SellerID.Int64)
//
//	// Way 2: update from variants
//	cacheModel, ok := util.GetOne[cache.Product](s.redis.Redis, payload.Id)
//	if !ok {
//		logger.Info("Product not in cache, ignore", "id", payload.Id)
//		return nil
//	}
//	if cacheModel.GetVersion() >= payload.GetVersion() {
//		logger.Info("Product have version greater than request, ignore", "cacheModel", cacheModel)
//		return nil
//	}
//
//	err := json.Unmarshal(payload.Variants, &cacheModel)
//	if err != nil {
//		logger.Error(err, "Fail unmarshal variants")
//		return err
//	}
//	// Update version
//	err = cacheModel.UpdateVersion(payload.GetVersion())
//	if err != nil {
//		logger.Error(err, "Update version fail") // ignore if fail
//	}
//
//	key, value := util.FuncConvertModel2Cache(payload.Id, cacheModel)
//	err = s.redis.Set(ctx, key, value)
//	if err != nil {
//		logger.Error(err, "Update to redis fail")
//		return err
//	}
//	logger.Info("Success")
//	return nil
//}

func (s *Service) UpdateRedisV2(ctx context.Context, ids []int64) error {
	logger := s.log.WithName("UpdateMemoryCache")
	logger.Info("Start process", "ids", ids)

	values, _ := util.GetMultiple[cache.Product](s.redis, ids)
	if len(values) == util.ZeroLength {
		return nil
	}

	products, err := s.storeDb.GetProductAndRelations(ctx, maps.Keys(values))
	if err != nil {
		logger.Error(err, "Call db get products fail")
	}
	if len(products) == util.ZeroLength {
		logger.Info("Not found products", "values", maps.Keys(values))
		return nil
	}
	// Convert list to map
	cacheModels := slice.Map(products, func(product store.GetProductAndRelationsRow) cache.Product {
		var cacheModel cache.Product
		cacheModel.FromDbV2(storeProduct.GetProductAndRelationsRow(product))
		return cacheModel
	})
	// Check version
	mapCacheModel := slice.KeyBy(cacheModels, func(c cache.Product) (int64, cache.Product) { return c.GetId(), c })
	oldObjects, _ := util.GetMultiple[cache.Product](s.redis, slice.Map(products, func(p store.GetProductAndRelationsRow) int64 { return p.ID }))
	needUpdateObject := make(map[string]any)
	for _, o := range oldObjects {
		newCache := mapCacheModel[o.GetId()]
		if newCache.GetVersion() > o.GetVersion() {
			k, v := util.FuncConvertModel2Cache(newCache)
			needUpdateObject[k] = v
		} else {
			delete(mapCacheModel, o.GetId())
		}
	}
	if len(needUpdateObject) == util.ZeroLength {
		return nil
	}

	err = s.redis.SetList(ctx, needUpdateObject)
	if err != nil {
		logger.Error(err, "Update to redis fail")
		return err
	}
	go func(cacheModels []cache.Product) {
		ctx := context.Background()
		t := time.Now()
		err = s.producer.Publish(ctx, util.TopicUpdateCache, producer.ProducerEvent{
			Key:   t.Format(time.RFC3339),
			Value: util.UpdateCacheEventValue{Objects: cacheModels},
		})
		if err != nil {
			s.log.Error(err, "Publish event fail")
			return
		}
	}(maps.Values(mapCacheModel))
	logger.Info("Success")
	return nil
}
