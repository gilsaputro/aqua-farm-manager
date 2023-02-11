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

func TestPondHandler_GetPondHandler(t *testing.T) {
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
			body: `{"size":2,"cursor":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetAllPond(gomock.Any(), gomock.Any()).Return(
					[]pond.GetPondInfoResponse{
						{
							ID:           1,
							Name:         "1",
							Capacity:     1,
							Depth:        1,
							WaterQuality: 1,
							Species:      "1",
							FarmID:       1,
						},
						{
							ID:           2,
							Name:         "2",
							Capacity:     2,
							Depth:        2,
							WaterQuality: 2,
							Species:      "2",
							FarmID:       2,
						},
					}, 2, nil,
				)
			},
			want: want{
				body: `{"data":{"ponds":[{"id":1,"name":"1","capacity":1,"depth":1,"water_quality":1,"species":"1","farm_id":1},{"id":2,"name":"2","capacity":2,"depth":2,"water_quality":2,"species":"2","farm_id":2}],"cursor":2},"code":200,"message":"success"}`,
				code: 200,
			},
		},
		{
			name: "error no data flow",
			body: `{"size":2,"cursor":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetAllPond(gomock.Any(), gomock.Any()).Return(
					[]pond.GetPondInfoResponse{}, 0, nil,
				)
			},
			want: want{
				body: `{"code":404,"message":"Data Not Found"}`,
				code: 404,
			},
		},
		{
			name: "timeout flow",
			body: `{"size":2,"cursor":1}`,
			args: args{
				timeout: 0,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetAllPond(gomock.Any(), gomock.Any()).Return(
					[]pond.GetPondInfoResponse{}, 0, nil,
				).AnyTimes()
			},
			want: want{
				body: `{"code":504,"message":"Timeout"}`,
				code: 504,
			},
		},
		{
			name: "data not found flow",
			body: `{"size":2,"cursor":1}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetAllPond(gomock.Any(), gomock.Any()).Return(
					[]pond.GetPondInfoResponse{}, 0, fmt.Errorf("record not found"),
				)
			},
			want: want{
				body: `{"code":404,"message":"Data Not Found"}`,
				code: 404,
			},
		},
		{
			name: "internal server flow",
			body: `{"size":0,"cursor":0}`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetAllPond(gomock.Any(), gomock.Any()).Return(
					[]pond.GetPondInfoResponse{}, 0, fmt.Errorf("some error"),
				)
			},
			want: want{
				body: `{"code":500,"message":"some error"}`,
				code: 500,
			},
		},
		{
			name: "broken request",
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

			r := httptest.NewRequest(http.MethodGet, "/pond", strings.NewReader(tt.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.GetPondHandler(w, r)
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
