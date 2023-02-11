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

	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
)

func TestFarmHandler_DeleteByIDFarmHandler(t *testing.T) {
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
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().DeleteFarmsWithDependencies(gomock.Any()).Return(farm.DeleteAllResponse{
					ID:      1,
					Name:    "a",
					PondIds: []uint{1, 2, 3},
				}, nil)
			},
			want: want{
				body: `{"data":{"id":1,"name":"a","ponds_id":[1,2,3]},"code":200,"message":"success"}`,
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
				farmDomain.EXPECT().DeleteFarmsWithDependencies(gomock.Any()).Return(farm.DeleteAllResponse{}, nil).AnyTimes()
			},
			want: want{
				body: `{"code":504,"message":"Timeout"}`,
				code: 504,
			},
		},
		{
			name: "error invalid farm request",
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().DeleteFarmsWithDependencies(gomock.Any()).Return(farm.DeleteAllResponse{}, farm.ErrInvalidFarm)
			},
			want: want{
				body: `{"code":404,"message":"Farm Is Not Exists"}`,
				code: 404,
			},
		},
		{
			name: "internal server error request",
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(farmDomain mock_farm.MockFarmDomain) {
				farmDomain.EXPECT().DeleteFarmsWithDependencies(gomock.Any()).Return(farm.DeleteAllResponse{}, fmt.Errorf("some error"))
			},
			want: want{
				body: `{"code":500,"message":"some error"}`,
				code: 500,
			},
		},
		{
			name: "invalid request",
			body: `a`,
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

			r := httptest.NewRequest(http.MethodDelete, "/pond/{id}", strings.NewReader(""))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)

			vars := map[string]string{
				"id": tt.body,
			}
			r = mux.SetURLVars(r, vars)

			w := httptest.NewRecorder()
			handler.DeleteByIDFarmHandler(w, r)
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
