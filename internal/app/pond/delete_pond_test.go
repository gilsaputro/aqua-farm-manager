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

func TestPondHandler_DeletePondHandler(t *testing.T) {
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
			body: `{"name":"a"}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().DeletePondInfo(gomock.Any()).Return(pond.DeleteDomainResponse{
					Name: "a",
					ID:   1,
				}, nil)
			},
			want: want{
				body: `{"data":{"id":1,"name":"a"},"code":200,"message":"success"}`,
				code: 200,
			},
		},
		{
			name: "timeout flow",
			body: `{"name":"a"}`,
			args: args{
				timeout: 0,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().DeletePondInfo(gomock.Any()).Return(pond.DeleteDomainResponse{
					Name: "a",
					ID:   1,
				}, nil).AnyTimes()
			},
			want: want{
				body: `{"code":504,"message":"Timeout"}`,
				code: 504,
			},
		},
		{
			name: "error invalid pond flow",
			body: `{"name":"a"}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().DeletePondInfo(gomock.Any()).Return(pond.DeleteDomainResponse{}, pond.ErrInvalidPond)
			},
			want: want{
				body: `{"code":404,"message":"Pond Is Not Exists"}`,
				code: 404,
			},
		},
		{
			name: "error internal server error flow",
			body: `{"name":"a"}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().DeletePondInfo(gomock.Any()).Return(pond.DeleteDomainResponse{}, fmt.Errorf("some error"))
			},
			want: want{
				body: `{"code":500,"message":"some error"}`,
				code: 500,
			},
		},
		{
			name: "error invalid request flow",
			body: `{"name":""}`,
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
			name: "error invalid duplicate request flow",
			body: `{
				"name": "a",
				"id": 1
			}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
			},
			want: want{
				body: `{"code":400,"message":"Please Choose to delete by ID or Name"}`,
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

			r := httptest.NewRequest(http.MethodDelete, "/pond", strings.NewReader(tt.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.DeletePondHandler(w, r)
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
