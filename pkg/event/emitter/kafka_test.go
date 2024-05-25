package emitter

import (
	"context"
	"errors"
	"testing"

	"github.com/douglasdgoulart/video-editor-api/internal/mocks"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/stretchr/testify/mock"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestKafkaEmitter_Send(t *testing.T) {
	tests := []struct {
		name    string
		event   event.Event
		wantErr error
		setup   func(m *mocks.KgoClientMock)
	}{
		{
			name:    "success",
			event:   event.Event{},
			wantErr: nil,
			setup: func(m *mocks.KgoClientMock) {
				m.On("ProduceSync", mock.Anything, mock.Anything).Return(kgo.ProduceResults{})
			},
		},
		{
			name:    "failure",
			event:   event.Event{},
			wantErr: errors.New("error"),
			setup: func(m *mocks.KgoClientMock) {
				m.On("ProduceSync", mock.Anything, mock.Anything).Return(kgo.ProduceResults{{Err: errors.New("error")}})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mocks.KgoClientMock)
			tt.setup(mockClient)

			k := &KafkaEmitter{
				cl: mockClient,
			}

			if err := k.Send(context.Background(), tt.event); (err != nil) != (tt.wantErr != nil) {
				t.Errorf("KafkaEmitter.Send() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
