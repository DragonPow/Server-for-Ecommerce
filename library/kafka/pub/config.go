package producer

import "time"

type producerConfig struct {
	maxNumberRetry    int
	maxPublishTimeOut time.Duration
	timeSleepPerRetry time.Duration
}

func loadDefaultConfig() *producerConfig {
	return &producerConfig{
		maxNumberRetry:    3,
		maxPublishTimeOut: 2 * time.Second,
		timeSleepPerRetry: 3 * time.Second,
	}
}

type ConfigOption = func(c *producerConfig)

func WithMaxNumberRetry(number int) ConfigOption {
	return func(c *producerConfig) {
		c.maxNumberRetry = number
	}
}

func WithTimeSleepPerRetry(t time.Duration) ConfigOption {
	return func(c *producerConfig) {
		c.timeSleepPerRetry = t
	}
}

func WithPublishTimeout(t time.Duration) ConfigOption {
	return func(c *producerConfig) {
		c.maxPublishTimeOut = t
	}
}
