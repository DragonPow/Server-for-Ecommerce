package service

import (
	"context"
	"encoding/json"
	producerDb "github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/producer"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/cache"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/redis_manager_service/util"
	"github.com/segmentio/kafka-go"
)

func (s *Service) Consume() error {
	kafkaConfig := s.cfg.KafkaConfig
	errChan := make(chan error)

	// Consumer update database
	updateConsumer := kafkaConfig.UpdateDbConsumer
	go func() {
		// create a new reader to the topic "update-db"
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     updateConsumer.Connections,
			Topic:       updateConsumer.Topic,
			GroupID:     updateConsumer.Group,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.LastOffset,
		})
		err := s.ProcessConsume(r, s.UpdateRedis)
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

func (s *Service) UpdateRedis(ctx context.Context, message kafka.Message) error {
	logger := s.log.WithName("UpdateMemoryCache").WithValues("message", message)
	logger.Info("Start process")
	var payload producerDb.UpdateDatabaseEventValue
	err := json.Unmarshal(message.Value, &payload)
	if err != nil {
		logger.Error(err, "Message value must be UpdateDatabaseEventValue")
		return err
	}

	// Way 1: Call db to enrich data
	//products, err := s.storeDb.GetProducts(ctx, []int64{payload.Id})
	//if err != nil {
	//	logger.Error(err, "Call db get products fail")
	//}
	//if len(products) == util.ZeroLength {
	//	err = fmt.Errorf("not found product with id %v", payload.Id)
	//	logger.Error(err, "Product not exists")
	//	return err
	//}
	//product := products[0]
	//templates, err := s.storeDb.GetProductTemplates(ctx, []int64{product.TemplateID.Int64})
	//if err != nil {
	//	logger.Error(err, "Call db get product templates fail")
	//}
	//if len(templates) == util.ZeroLength {
	//	err = fmt.Errorf("not found product template with id %v", payload.Id)
	//	logger.Error(err, "Product not exists")
	//	return err
	//}
	//template := templates[0]
	//var cacheModel cache.Product
	//cacheModel.FromDb(storeProduct.Product(product), template.CategoryID.Int64, template.UomID.Int64, template.SellerID.Int64)

	// Way 2: update from variants
	cacheModel, ok := util.GetOne[cache.Product](s.redis, payload.Id)
	if !ok {
		logger.Info("Product not in cache, ignore", "id", payload.Id)
		return nil
	}
	if cacheModel.GetVersion() >= payload.GetVersion() {
		logger.Info("Product have version greater than request, ignore", "cacheModel", cacheModel)
		return nil
	}

	err = json.Unmarshal(payload.Variants, &cacheModel)
	if err != nil {
		logger.Error(err, "Fail unmarshal variants")
		return err
	}
	// Update version
	err = cacheModel.UpdateVersion(payload.GetVersion())
	if err != nil {
		logger.Error(err, "Update version fail") // ignore if fail
	}

	key, value := util.FuncConvertModel2Cache(payload.Id, cacheModel)
	err = s.redis.Set(ctx, key, value)
	if err != nil {
		logger.Error(err, "Update to redis fail")
		return err
	}
	logger.Info("Success")
	return nil
}
