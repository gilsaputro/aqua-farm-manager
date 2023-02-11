package trackingevent

import (
	"aqua-farm-manager/internal/domain/stat"
	"aqua-farm-manager/internal/domain/stat/mock_stat"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nsqio/go-nsq"
)

func TestNewTrackingEverntConsumer(t *testing.T) {
	type args struct {
		topic        string
		channel      string
		host         string
		maxInFlight  int
		numconsumer  int
		timeoutInSec int
		stat         stat.StatDomain
	}
	tests := []struct {
		name string
		args args
		want *TrackingEventConsumer
	}{
		{
			name: "success",
			args: args{
				topic:        "topic",
				channel:      "channel",
				host:         "host",
				maxInFlight:  1,
				numconsumer:  1,
				timeoutInSec: 1,
				stat:         &stat.Stat{},
			},
			want: &TrackingEventConsumer{
				topic:        "topic",
				channel:      "channel",
				host:         "host",
				maxInFlight:  1,
				numconsumer:  1,
				timeoutInSec: 1,
				stat:         &stat.Stat{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrackingEventConsumer(tt.args.topic, tt.args.channel, tt.args.host, tt.args.maxInFlight, tt.args.numconsumer, tt.args.timeoutInSec, tt.args.stat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrackingEventConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockNSQ struct {
	nsq.Message
}

func (m *mockNSQ) OnFinish(msg *nsq.Message) {}

func (m *mockNSQ) OnRequeue(msg *nsq.Message, delay time.Duration, backoff bool) {}

func (m *mockNSQ) OnTouch(*nsq.Message) {}

func TestTrackingEventConsumer_HandleMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	domain := mock_stat.NewMockStatDomain(mockCtrl)
	type args struct {
		body string
	}
	tests := []struct {
		name     string
		body     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success flow",
			body: `{
				"path": "/v1/farms",
				"code": 200,
				"method": "GET",
				"ua": "Mozilla/5.0"
			  }`,
			mockFunc: func() {
				domain.EXPECT().IngestStatAPI(gomock.Any())
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			body: `{
				"path": "",
				"code": 200,
				"method": "GET",
				"ua": "Mozilla/5.0"
			  }`,
			mockFunc: func() {
			},
			wantErr: false,
		},
		{
			name: "error unmarshall",
			body: `{
				"path": "",
				"code": 200,
				"method": "GET",
				"ua": "Mozilla/5.0",
			  }`,
			mockFunc: func() {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			c := TrackingEventConsumer{
				topic:        "topic",
				channel:      "channel",
				host:         "host",
				maxInFlight:  1,
				numconsumer:  1,
				timeoutInSec: 1,
				stat:         domain,
			}
			msg := &nsq.Message{
				Body:     []byte(tt.body),
				Delegate: &mockNSQ{},
			}
			if err := c.HandleMessage(msg); (err != nil) != tt.wantErr {
				t.Errorf("TrackingEventConsumer.HandleMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
