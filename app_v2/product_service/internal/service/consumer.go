package service

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func (s *Service) Consume() error {
	kafkaConfig := s.cfg.KafkaConfig
	// create a new reader to the topic "my-topic"
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaConfig.Connections,
		Topic:    kafkaConfig.Topic,
		GroupID:  kafkaConfig.ConsumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	// consume messages
	for {
		ctx := context.Background()
		m, err := r.ReadMessage(ctx)
		if err != nil {
			s.log.Error(err, "Read message fail")
			return err
		}
		s.log.Info("receive message offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		go func() {
			err := s.UpdateMemoryCache()
			if err != nil {
				s.log.Error(err, "UpdateMemoryCache fail")
			}
		}()
	}
}

func (s *Service) UpdateMemoryCache() error {
	return nil
}
