package service

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/app_v2/product_service/util"
	pubEvent "Server-for-Ecommerce/app_v2/redis_manager_service/util"
	"Server-for-Ecommerce/library/slice"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io"
	"time"
)

func (s *Service) Consume() error {
	kafkaConfig := s.cfg.KafkaConfig
	errChan := make(chan error)

	// Consumer update database
	updateConsumer := kafkaConfig.UpdateDbConsumer
	go func() {
		// create a new reader to the topic "my-topic"
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     kafkaConfig.Connections,
			Topic:       updateConsumer.Topic,
			GroupID:     fmt.Sprintf(updateConsumer.Group + time.Now().String()),
			StartOffset: kafka.LastOffset,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			//StartOffset:            kafka.LastOffset,
			MaxWait:                1 * time.Second,
			PartitionWatchInterval: 1 * time.Second,
		})
		//err := s.ProcessConsume(r, s.UpdateMemoryCache)
		//if err != nil {
		//	errChan <- err
		//}
		err := s.ProcessConsume(r, s.UpdateMemoryCacheV2)
		if err != nil {
			errChan <- err
		}
	}()

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
				s.log.Info("Message EOF")
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

//func (s *Service) UpdateMemoryCache(ctx context.Context, message kafka.Message) error {
//	logger := s.log.WithName("UpdateMemoryCache").WithValues("message", message)
//	logger.Info("Start process")
//	var payload producerDb.UpdateDatabaseEventValue
//	err := json.Unmarshal(message.Value, &payload)
//	if err != nil {
//		logger.Error(err, "Message value must be UpdateDatabaseEventValue")
//		return err
//	}
//	//products, err := s.storeDb.GetProducts(ctx, []int64{payload.Id})
//	//if err != nil {
//	//	logger.Error(err, "Call db get products fail")
//	//}
//	//if len(products) == util.ZeroLength {
//	//	err = fmt.Errorf("not found product with id %v", payload.Id)
//	//	logger.Error(err, "Product not exists")
//	//	return err
//	//}
//
//	err = s.memCache.SetProductByAttr(cache.Product{ID: payload.Id}, payload.Variants, payload.GetVersion())
//	if err != nil {
//		logger.Error(err, "Fail to SetProductByAttr")
//		return err
//	}
//	logger.Info("Success")
//	return nil
//}

func (s *Service) UpdateMemoryCacheV2(ctx context.Context, message kafka.Message) error {
	logger := s.log.WithName("UpdateMemoryCacheV2").WithValues("message", message)
	logger.Info("Start process")
	var payload pubEvent.UpdateCacheEventValue
	err := json.Unmarshal(message.Value, &payload)
	if err != nil {
		logger.Error(err, "Message value must be UpdateCacheEventValue")
		return err
	}

	oldObjects, _ := s.memCache.GetListProduct(slice.Map(payload.Objects, func(o cache.Product) int64 {
		return o.GetId()
	}))
	mapCacheModel := slice.KeyBy(payload.Objects, func(o cache.Product) (int64, cache.ModelValue) {
		return o.GetId(), o
	})
	needUpdateObject := make(map[int64]cache.ModelValue)
	for _, o := range oldObjects {
		newCache := mapCacheModel[o.GetId()]
		if newCache.GetVersion() > o.GetVersion() {
			needUpdateObject[newCache.GetId()] = newCache
		}
	}
	if len(needUpdateObject) == util.ZeroLength {
		return nil
	}

	err = s.memCache.SetMultiple(needUpdateObject)
	if err != nil {
		logger.Error(err, "Fail to SetProductByAttr")
		return err
	}
	logger.Info("Success")
	return nil
}
