package receiver

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/twmb/franz-go/pkg/kfake"
	"github.com/twmb/franz-go/pkg/kgo"
)

type MockHandleFunc struct {
	mock.Mock
}

func (m *MockHandleFunc) Handle(event *event.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestKafkaEventReceiver_Receive(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name      string
		args      args
		setup     func(m *MockHandleFunc) *kfake.Cluster
		assertion func(m *MockHandleFunc)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			setup: func(m *MockHandleFunc) *kfake.Cluster {
				m.On("Handle", mock.Anything).Return(nil)

				fakeKafka, _ := kfake.NewCluster(
					kfake.SeedTopics(-1, "topic"),
				)
				return fakeKafka
			},
			assertion: func(m *MockHandleFunc) {
				assert.Eventually(t, func() bool {
					return m.AssertExpectations(t)
				}, 5*time.Second, 1*time.Second)
			},
		},
		{
			name: "unmarshal failure",
			args: args{
				ctx: context.Background(),
			},
			setup: func(m *MockHandleFunc) *kfake.Cluster {
				m.On("Handle", mock.Anything).Return(nil)

				fakeKafka, _ := kfake.NewCluster(
					kfake.SeedTopics(-1, "topic"),
				)
				return fakeKafka
			},
			assertion: func(m *MockHandleFunc) {
				time.Sleep(2 * time.Second)
				m.AssertNotCalled(t, "Handle")
			},
		},
		{
			name: "handle failure",
			args: args{
				ctx: context.Background(),
			},
			setup: func(m *MockHandleFunc) *kfake.Cluster {
				m.On("Handle", mock.Anything).Return(nil)

				fakeKafka, _ := kfake.NewCluster(
					kfake.SeedTopics(-1, "topic"),
				)
				return fakeKafka
			},
			assertion: func(m *MockHandleFunc) {
				assert.Eventually(t, func() bool {
					return len(m.Calls) == 5
				}, 5*time.Second, 1*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerMock := new(MockHandleFunc)
			fakeKafka := tt.setup(handlerMock)

			defer fakeKafka.Close()

			cl, err := kgo.NewClient(
				kgo.SeedBrokers(fakeKafka.ListenAddrs()...),
				kgo.ConsumeTopics("topic"),
				kgo.ConsumerGroup("group"),
			)
			assert.NoError(t, err)

			k := &KafkaEventReceiver{
				cl:     cl,
				logger: slog.Default(),
			}

			go func() {
				if tt.name == "unmarshal failure" {
					produceInvalidMessage(cl, "topic")
				} else {
					produceValidMessage(cl, "topic")
				}
			}()

			go k.Receive(tt.args.ctx, handlerMock.Handle)

			assert.Eventually(t, func() bool {
				if tt.name == "unmarshal failure" {
					return handlerMock.AssertNotCalled(t, "Handle")
				}
				return handlerMock.AssertExpectations(t)
			}, 5*time.Second, 1*time.Second)
		})
	}
}

func produceValidMessage(fakeKafka *kgo.Client, topic string) {
	e := event.Event{}
	value, _ := json.Marshal(e)

	fakeKafka.ProduceSync(context.Background(), &kgo.Record{
		Topic: topic,
		Value: value,
	})
}

func produceInvalidMessage(fakeKafka *kgo.Client, topic string) {
	fakeKafka.ProduceSync(context.Background(), &kgo.Record{
		Topic: topic,
		Value: []byte("invalid json"),
	})
}
