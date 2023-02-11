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
	"github.com/gorilla/mux"
)

func TestFarmHandler_GetByIDFarmHandler(t *testing.T) {
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
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarmInfoByID(gomock.Any()).Return(farm.GetFarmInfoResponse{
					ID:       1,
					Name:     "name",
					Location: "loc",
					Owner:    "own",
					Area:     "area",
					PondInfos: []farm.PondInfo{
						{
							ID:           1,
							Name:         "p1",
							Capacity:     1,
							Depth:        1,
							WaterQuality: 1,
							Species:      "1",
						},
					},
				}, nil)
			},
			want: want{
				body: `{"data":{"id":1,"name":"name","location":"loc","owner":"own","area":"area","pond_info":[{"id":1,"name":"p1","capacity":1,"depth":1,"water_quality":1,"species":"1","status":0}]},"code":200,"message":"success"}`,
				code: 200,
			},
		},
		{
			name: "timeout flow",
			body: `1`,
			args: args{
				timeout: 0,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarmInfoByID(gomock.Any()).Return(farm.GetFarmInfoResponse{}, nil).AnyTimes()
			},
			want: want{
				body: `{"code":504,"message":"Timeout"}`,
				code: 504,
			},
		},
		{
			name: "invalid request flow",
			body: `A`,
			args: args{
				timeout: 0,
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
			name: "error data not found flow",
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarmInfoByID(gomock.Any()).Return(farm.GetFarmInfoResponse{}, fmt.Errorf("record not found"))
			},
			want: want{
				body: `{"code":404,"message":"Data Not Found"}`,
				code: 404,
			},
		},
		{
			name: "internal error flow",
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().GetFarmInfoByID(gomock.Any()).Return(farm.GetFarmInfoResponse{}, fmt.Errorf("some error"))
			},
			want: want{
				body: `{"code":500,"message":"some error"}`,
				code: 500,
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

			r := httptest.NewRequest(http.MethodPost, "/farm/{id}", strings.NewReader(""))

			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)

			vars := map[string]string{
				"id": tt.body,
			}
			r = mux.SetURLVars(r, vars)

			w := httptest.NewRecorder()
			handler.GetByIDFarmHandler(w, r)
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
