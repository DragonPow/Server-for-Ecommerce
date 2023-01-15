package producer

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

func funcConvertEvent2Message(topicName string) func(e ProducerEvent) kafka.Message {
	return func(e ProducerEvent) kafka.Message {
		v, err := json.Marshal(e.Value)
		if err != nil {
			panic("marshal fail")
		}
		return kafka.Message{
			Topic: topicName,
			Key:   []byte(e.Key),
			Value: v,
			Time:  time.Now(),
		}
	}
}
