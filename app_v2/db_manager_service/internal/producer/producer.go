package producer

import (
	"context"
	"errors"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/config"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
	"github.com/go-logr/logr"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer interface {
	Publish(ctx context.Context, topicName string, events ...ProducerEvent) error
}

type producer struct {
	topics            map[string]*kafka.Writer
	log               logr.Logger
	maxNumberRetry    int
	maxPublishTimeOut time.Duration
	connections       []string
	timeSleepPerRetry time.Duration
}

func NewProducer(cfg config.KafkaConfig, log logr.Logger) (Producer, error) {
	result := &producer{
		log:               log,
		connections:       cfg.Connections,
		topics:            make(map[string]*kafka.Writer, len(cfg.ListProducer)),
		maxNumberRetry:    cfg.MaxNumberRetry,
		timeSleepPerRetry: time.Duration(cfg.TimeSleepPerRetryMillisecond) * time.Millisecond,
		maxPublishTimeOut: time.Duration(cfg.MaxPublishTimeoutSecond) * time.Second,
	}
	// Register topic
	for _, producer := range cfg.ListProducer {
		err := result.RegisterTopic(producer)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (p *producer) RegisterTopic(producer config.Producer) error {
	if _, ok := p.topics[producer.Topic]; ok {
		return fmt.Errorf("topic %s is duplicate", producer.Topic)
	}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(p.connections...),
		Topic:                  producer.Topic,
		Balancer:               &kafka.RoundRobin{},
		WriteTimeout:           p.maxPublishTimeOut,
		Async:                  true,
		Completion:             nil,
		AllowAutoTopicCreation: true,
	}
	p.topics[producer.Topic] = writer
	return nil
}

func (p *producer) Publish(ctx context.Context, topicName string, events ...ProducerEvent) error {
	logger := p.log.WithName("Publish").WithValues("topicName", topicName)
	topic, ok := p.topics[topicName]
	if !ok {
		logger.Error(fmt.Errorf("topic not register"), "fail to get message")
		return fmt.Errorf("topic not register")
	}
	if len(events) == util.ZeroLength {
		return nil
	}

	for i := 0; i < p.maxNumberRetry; i++ {
		newCtx, cancel := context.WithTimeout(ctx, p.maxPublishTimeOut)
		defer cancel()

		// Push Message
		err := topic.WriteMessages(newCtx, math.Convert(events, funcConvertEvent2Message(topicName))...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			logger.Error(err, "Push message fail, sleep and retry")
			time.Sleep(p.timeSleepPerRetry)
			continue
		}
		if err != nil {
			logger.Error(err, "Push message fail")
			return err
		}
		topic.Close()
	}
	logger.Info("Push message success")
	return nil
}
