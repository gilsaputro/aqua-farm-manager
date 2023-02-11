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

func TestFarmHandler_GetFarmHandler(t *testing.T) {
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
		mockFunc    func(farmDomain mock_farm.MockFarmDomain)
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
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarm(gomock.Any(), gomock.Any()).Return(
					[]farm.GetFarmInfoResponse{
						{
							ID:       1,
							Name:     "1",
							Location: "1",
							Owner:    "1",
							Area:     "1",
							PondIDs:  []uint{1, 2, 3},
						},
						{
							ID:       2,
							Name:     "2",
							Location: "2",
							Owner:    "2",
							Area:     "2",
							PondIDs:  []uint{4, 5, 6},
						},
					}, 2, nil,
				)
			},
			want: want{
				body: `{"data":{"farms":[{"id":1,"name":"1","location":"1","owner":"1","area":"1","list_pondID":[1,2,3]},{"id":2,"name":"2","location":"2","owner":"2","area":"2","list_pondID":[4,5,6]}],"cursor":2},"code":200,"message":"success"}`,
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
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarm(gomock.Any(), gomock.Any()).Return(
					[]farm.GetFarmInfoResponse{}, 0, nil,
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
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarm(gomock.Any(), gomock.Any()).Return(
					[]farm.GetFarmInfoResponse{}, 0, nil,
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
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarm(gomock.Any(), gomock.Any()).Return(
					[]farm.GetFarmInfoResponse{}, 0, fmt.Errorf("record not found"),
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
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarm(gomock.Any(), gomock.Any()).Return(
					[]farm.GetFarmInfoResponse{}, 0, fmt.Errorf("some error"),
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

			r := httptest.NewRequest(http.MethodGet, "/farm", strings.NewReader(tt.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.GetFarmHandler(w, r)
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
