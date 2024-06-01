package emitter

import (
	"context"
	"encoding/json"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaEmitter struct {
	cl    event.KgoClient
	topic string
}

func NewKafkaEmitter(cfg *configuration.KafkaProducerConfig) EventEmitter {
	cl, err := kgo.NewClient(kgo.SeedBrokers(cfg.Brokers...))
	if err != nil {
		panic(err)
	}
	return &KafkaEmitter{
		cl:    cl,
		topic: cfg.Topic,
	}
}

func (k *KafkaEmitter) Send(ctx context.Context, event event.Event) error {
	serializedEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = k.cl.ProduceSync(ctx, &kgo.Record{
		Topic: k.topic,
		Value: serializedEvent,
		Key:   []byte(event.Id),
	}).FirstErr()

	return err
}
