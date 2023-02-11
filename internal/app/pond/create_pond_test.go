package pond

import (
	"aqua-farm-manager/internal/domain/pond"
	"aqua-farm-manager/internal/domain/pond/mock_pond"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPondHandler_CreatePondHandler(t *testing.T) {
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
		mockFunc    func(pondDomain mock_pond.MockPondDomain)
		mockContext func() (context.Context, func())
		want        want
	}{
		{
			name: "success flow",
			body: `{"name":"Pond 1","capacity":1000,"depth":2.5,"water_quality":7.8,"species":"Tilapia","farm_id":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().CreatePondInfo(gomock.Any()).Return(pond.CreateDomainResponse{
					PondID: 1,
				}, nil)
			},
			want: want{
				body: `{"data":{"pond_id":1},"code":200,"message":"success"}`,
				code: 200,
			},
		},
		{
			name: "timeout flow",
			body: `{"name":"Pond 1","capacity":1000,"depth":2.5,"water_quality":7.8,"species":"Tilapia","farm_id":1}`,
			args: args{
				timeout: 0,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().CreatePondInfo(gomock.Any()).Return(pond.CreateDomainResponse{
					PondID: 1,
				}, nil).AnyTimes()
			},
			want: want{
				body: `{"code":504,"message":"Timeout"}`,
				code: 504,
			},
		},
		{
			name: "error duplicate pond flow",
			body: `{"name":"Pond 1","capacity":1000,"depth":2.5,"water_quality":7.8,"species":"Tilapia","farm_id":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().CreatePondInfo(gomock.Any()).Return(pond.CreateDomainResponse{}, pond.ErrDuplicatePond)
			},
			want: want{
				body: `{"code":409,"message":"Pond Is Already Exists"}`,
				code: 409,
			},
		},
		{
			name: "error invalid farm flow",
			body: `{"name":"Pond 1","capacity":1000,"depth":2.5,"water_quality":7.8,"species":"Tilapia","farm_id":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().CreatePondInfo(gomock.Any()).Return(pond.CreateDomainResponse{}, pond.ErrInvalidFarm)
			},
			want: want{
				body: `{"code":404,"message":"Farm Is Not Exists"}`,
				code: 404,
			},
		},
		{
			name: "error invalid farm flow",
			body: `{"name":"Pond 1","capacity":1000,"depth":2.5,"water_quality":7.8,"species":"Tilapia","farm_id":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().CreatePondInfo(gomock.Any()).Return(pond.CreateDomainResponse{}, fmt.Errorf("Internal Server Error"))
			},
			want: want{
				body: `{"code":500,"message":"Internal Server Error"}`,
				code: 500,
			},
		},
		{
			name: "error invalid request flow",
			body: `{"name":"","capacity":1000,"depth":2.5,"water_quality":7.8,"species":"Tilapia","farm_id":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
			},
			want: want{
				body: `{"code":400,"message":"Invalid Parameter Request"}`,
				code: 400,
			},
		},
		{
			name: "error broken request flow",
			body: `{`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
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
			pondDomain := mock_pond.NewMockPondDomain(mockCtrl)
			tt.mockFunc(*pondDomain)

			handler := PondHandler{
				domain:       pondDomain,
				timeoutInSec: tt.args.timeout,
			}

			r := httptest.NewRequest(http.MethodPost, "/pond", strings.NewReader(tt.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.CreatePondHandler(w, r)
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
