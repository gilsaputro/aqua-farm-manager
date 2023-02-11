package middleware

import (
	"aqua-farm-manager/pkg/nsq"
	"aqua-farm-manager/pkg/nsq/mock_nsq"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewMiddleware(t *testing.T) {
	type args struct {
		topic string
		nsq   nsq.NsqMethod
	}
	tests := []struct {
		name string
		args args
		want Middleware
	}{
		{
			name: "success",
			args: args{
				topic: "topic",
				nsq:   &nsq.Client{},
			},
			want: Middleware{
				nsq:   &nsq.Client{},
				topic: "topic",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMiddleware(tt.args.topic, tt.args.nsq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiddleware_Middleware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nsqMock := mock_nsq.NewMockNsqMethod(mockCtrl)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	nsqMock.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	tests := []struct {
		name       string
		status     int
		ua         string
		path       string
		method     string
		want       int
		wantUa     string
		wantPath   string
		wantMethod string
	}{
		{
			name:       "Success",
			status:     http.StatusOK,
			ua:         "ua",
			path:       "/v1/test",
			method:     "GET",
			want:       http.StatusOK,
			wantUa:     "ua",
			wantPath:   "/v1/test",
			wantMethod: "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Middleware{
				nsq: nsqMock,
			}

			var gotPath string
			var gotMethod string
			var gotUa string
			var gotStatus int

			middleware := m.Middleware(next)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, tt.path, nil)
			request.Header.Set("User-Agent", tt.ua)

			middleware(recorder, request)

			gotPath = request.URL.Path
			gotMethod = request.Method
			gotUa = request.Header.Get("User-Agent")
			gotStatus = recorder.Result().StatusCode

			if gotPath != tt.wantPath {
				t.Errorf("Middleware.Middleware() Path = %v, want %v", gotPath, tt.wantPath)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("Middleware.Middleware() Method = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotUa != tt.wantUa {
				t.Errorf("Middleware.Middleware() User Agent = %v, want %v", gotUa, tt.wantUa)
			}
			if gotStatus != tt.want {
				t.Errorf("Middleware.Middleware() Status = %v, want %v", gotStatus, tt.want)
			}
		})
	}
}
