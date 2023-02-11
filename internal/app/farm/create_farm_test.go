package farm

import (
	"aqua-farm-manager/internal/domain/farm"
	"aqua-farm-manager/internal/domain/farm/mock_farm"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestFarmHandler_CreateFarmHandler(t *testing.T) {
	type args struct {
		timeout int
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name        string
		body        string
		args        args
		mockFunc    func(pondDomain mock_farm.MockFarmDomain)
		mockContext func() (context.Context, func())
		want        want
	}{
		{
			name: "success flow",
			body: `{ "name": "Green Pastures 2", "location": "California", "owner": "Jane Doe", "area": "7.5 Acres" }`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().CreateFarmInfo(gomock.Any()).Return(farm.CreateDomainResponse{
					ID: 1,
				}, nil)
			},
			want: want{
				body: `{"data":{"id":1},"code":200,"message":"success"}`,
				code: 200,
			},
		},
		{
			name: "timeout flow",
			body: `{ "name": "Green Pastures 2", "location": "California", "owner": "Jane Doe", "area": "7.5 Acres" }`,
			args: args{
				timeout: 0,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().CreateFarmInfo(gomock.Any()).Return(farm.CreateDomainResponse{}, nil).AnyTimes()
			},
			want: want{
				body: `{"code":504,"message":"Timeout"}`,
				code: 504,
			},
		},
		{
			name: "internal server error flow",
			body: `{ "name": "Green Pastures 2", "location": "California", "owner": "Jane Doe", "area": "7.5 Acres" }`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().CreateFarmInfo(gomock.Any()).Return(farm.CreateDomainResponse{}, fmt.Errorf("some error"))
			},
			want: want{
				body: `{"code":500,"message":"some error"}`,
				code: 500,
			},
		},
		{
			name: "invalid request flow",
			body: `{ "name": "", "location": "California", "owner": "Jane Doe", "area": "7.5 Acres" }`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
			},
			want: want{
				body: `{"code":400,"message":"Invalid Parameter Request"}`,
				code: 400,
			},
		},
		{
			name: "broken request flow",
			body: `{ `,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
			},
			want: want{
				body: `{"code":400,"message":"Bad Request"}`,
				code: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			farmDomain := mock_farm.NewMockFarmDomain(mockCtrl)
			tt.mockFunc(*farmDomain)

			handler := FarmHandler{
				domain:       farmDomain,
				timeoutInSec: tt.args.timeout,
			}

			r := httptest.NewRequest(http.MethodPost, "/farm", strings.NewReader(tt.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.CreateFarmHandler(w, r)
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
