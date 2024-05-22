package receiver

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaEventReceiver struct {
	cl event.KgoClient
}

func NewKafkaEventReceiver(cfg configuration.KafkaConsumerConfig) EventReceiver {
	getKgoOffset(cfg)

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topic),
	)
	if err != nil {
		panic(err)
	}
	return &KafkaEventReceiver{
		cl: cl,
	}
}

func getKgoOffset(cfg configuration.KafkaConsumerConfig) kgo.Offset {
	if cfg.Offset == "earliest" {
		return kgo.NewOffset().AtStart()
	}

	return kgo.NewOffset().AtEnd()
}

func (k *KafkaEventReceiver) Receive(ctx context.Context, handle func(event *event.Event) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fetches := k.cl.PollFetches(ctx)
			iter := fetches.RecordIter()
			for !iter.Done() {
				var e event.Event
				record := iter.Next()

				if err := json.Unmarshal(record.Value, &e); err != nil {
					slog.Error("error unmarshalling event", "error", err, "event", string(record.Value))
					continue
				}

				var processErr error
				for retryCount := range 5 {
					if processErr = handle(&e); processErr != nil {
						slog.Error("error handling event", "error", processErr, "retryCount", retryCount)
						continue
					}
					break
				}

				if processErr != nil {
					slog.Error("error handling event", "error", processErr)
				}
			}
		}
	}
}
