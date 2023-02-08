package stat

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"aqua-farm-manager/internal/domain/stat"
	"aqua-farm-manager/internal/domain/stat/mock_stat"

	"github.com/golang/mock/gomock"
)

func TestNewStatHandler(t *testing.T) {
	type args struct {
		stat    stat.StatDomain
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *StatHandler
	}{
		{
			name: "success without option",
			args: args{
				stat:    &stat.Stat{},
				options: []Option{},
			},
			want: &StatHandler{
				stat:         &stat.Stat{},
				timeoutInSec: 5,
			},
		},
		{
			name: "success wit option",
			args: args{
				stat:    &stat.Stat{},
				options: []Option{WithTimeoutOptions(10)},
			},
			want: &StatHandler{
				stat:         &stat.Stat{},
				timeoutInSec: 10,
			},
		},
		{
			name: "success wit invalid option value",
			args: args{
				stat:    &stat.Stat{},
				options: []Option{WithTimeoutOptions(0)},
			},
			want: &StatHandler{
				stat:         &stat.Stat{},
				timeoutInSec: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatHandler(tt.args.stat, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatHandler_GetStatHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		timeout int
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name        string
		args        args
		mockFunc    func(*mock_stat.MockStatDomain)
		mockContext func() (context.Context, func())
		want        want
	}{
		{
			name: "success flow",
			args: args{
				timeout: 5,
			},
			mockFunc: func(msd *mock_stat.MockStatDomain) {
				msd.EXPECT().GenerateStatAPI().Return(map[string]stat.StatMetrics{"POST /farms": {NumRequested: 3, NumUniqAgent: 1, NumSuccess: 2, NumError: 1}})
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 200,
				body: `{"data":{"POST /farms":{"count":3,"unique_user_agent":1,"num_success":2,"num_error":1}},"code":200,"message":"success"}`,
			},
		},
		{
			name: "got timeout flow",
			args: args{
				timeout: 0,
			},
			mockFunc: func(msd *mock_stat.MockStatDomain) {
				msd.EXPECT().GenerateStatAPI().Return(map[string]stat.StatMetrics{}).AnyTimes()
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 504,
				body: `{"code":504,"message":"Timeout"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := mock_stat.NewMockStatDomain(mockCtrl)
			tt.mockFunc(domain)

			handler := StatHandler{
				stat:         domain,
				timeoutInSec: tt.args.timeout,
			}

			r := httptest.NewRequest(http.MethodGet, "/stat", strings.NewReader(""))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.GetStatHandler(w, r)
			result := w.Result()
			resBody, err := ioutil.ReadAll(result.Body)

			if err != nil {
				t.Fatalf("Error read body err = %v\n", err)
			}

			if result.StatusCode != tt.want.code {
				t.Fatalf("GetStatHandler status code got =%d, want %d \n", result.StatusCode, tt.want.code)
			}

			if string(resBody) != tt.want.body {
				t.Fatalf("GetStatHandler body got =%s, want %s \n", string(resBody), tt.want.body)
			}
		})
	}
}
