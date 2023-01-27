package service

import (
	"context"
	"encoding/json"
	producerDb "github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/producer"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/cache"
	"github.com/segmentio/kafka-go"
)

func (s *Service) Consume() error {
	kafkaConfig := s.cfg.KafkaConfig
	errChan := make(chan error)

	// Consumer update database
	updateConsumer := kafkaConfig.UpdateDbConsumer
	go func() {
		// create a new reader to the topic "my-topic"
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     updateConsumer.Connections,
			Topic:       updateConsumer.Topic,
			GroupID:     updateConsumer.Group,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.LastOffset,
		})
		err := s.ProcessConsume(r, s.UpdateMemoryCache)
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

func (s *Service) UpdateMemoryCache(ctx context.Context, message kafka.Message) error {
	logger := s.log.WithName("UpdateMemoryCache").WithValues("message", message)
	logger.Info("Start process")
	var payload producerDb.UpdateDatabaseEventValue
	err := json.Unmarshal(message.Value, &payload)
	if err != nil {
		logger.Error(err, "Message value must be UpdateDatabaseEventValue")
		return err
	}
	//products, err := s.storeDb.GetProducts(ctx, []int64{payload.Id})
	//if err != nil {
	//	logger.Error(err, "Call db get products fail")
	//}
	//if len(products) == util.ZeroLength {
	//	err = fmt.Errorf("not found product with id %v", payload.Id)
	//	logger.Error(err, "Product not exists")
	//	return err
	//}

	err = s.memCache.SetProductByAttr(cache.Product{ID: payload.Id}, payload.Variants, payload.GetVersion())
	if err != nil {
		logger.Error(err, "Fail to SetProductByAttr")
		return err
	}
	logger.Info("Success")
	return nil
}
