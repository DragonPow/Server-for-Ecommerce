package producer

import (
	"Server-for-Ecommerce/app_v2/db_manager_service/config"
	"Server-for-Ecommerce/app_v2/db_manager_service/util"
	"Server-for-Ecommerce/library/math"
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer interface {
	Publish(ctx context.Context, topicName string, events ...ProducerEvent) error
	Register(topicName string) error
}

type producer struct {
	topics      map[string]*kafka.Writer
	log         *logr.Logger
	cfg         *producerConfig
	connections []string
}

func NewProducer(connections []string, log *logr.Logger, options ...ConfigOption) (p Producer, err error) {
	cfg := loadDefaultConfig()
	for _, opt := range options {
		opt(cfg)
	}
	result := &producer{
		log:         log,
		connections: connections,
		topics:      make(map[string]*kafka.Writer),
		cfg:         cfg,
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
		WriteTimeout:           p.cfg.maxPublishTimeOut,
		Async:                  true,
		Completion:             nil,
		AllowAutoTopicCreation: true,
	}
	p.topics[producer.Topic] = writer
	return nil
}

func (p *producer) Register(topicName string) error {
	if _, ok := p.topics[topicName]; ok {
		return fmt.Errorf("topic %s is exists", topicName)
	}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(p.connections...),
		Topic:                  topicName,
		Balancer:               &kafka.RoundRobin{},
		WriteTimeout:           p.cfg.maxPublishTimeOut,
		Async:                  true,
		Completion:             nil,
		AllowAutoTopicCreation: true,
	}
	p.topics[topicName] = writer
	return nil
}

func (p *producer) Publish(ctx context.Context, topicName string, events ...ProducerEvent) error {
	logger := p.log.WithName("Publish").WithValues("topicName", topicName)
	// Check topic exists register
	topic, ok := p.topics[topicName]
	if !ok {
		logger.Error(fmt.Errorf("topic not register"), "fail to get message")
		return fmt.Errorf("topic not register")
	}
	if len(events) == util.ZeroLength {
		return nil
	}

	//defer topic.Close()
	var retry int
	// Begin publish with retry
	for retry = 0; retry < p.cfg.maxNumberRetry; retry++ {
		newCtx, cancel := context.WithTimeout(ctx, p.cfg.maxPublishTimeOut)
		defer cancel()

		// Push Message
		messages := math.Convert(events, funcConvertEvent2Message(topicName))
		err := topic.WriteMessages(newCtx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			logger.Error(err, "Push message fail, sleep and retry")
			time.Sleep(p.cfg.timeSleepPerRetry)
			continue
		}
		if err != nil {
			logger.Error(err, "Push message fail")
			return err
		}
		break
	}
	if retry == p.cfg.maxNumberRetry {
		logger.Info("Push message fail, get max number retry")
		return nil
	}
	logger.Info("Push message success")
	return nil
}
