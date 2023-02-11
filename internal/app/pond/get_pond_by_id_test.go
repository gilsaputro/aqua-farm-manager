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
	"github.com/gorilla/mux"
)

func TestPondHandler_GetByIDPondHandler(t *testing.T) {
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
			body: `1`,
			args: args{
				timeout: 10,
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetPondInfoByID(gomock.Any()).Return(pond.GetPondInfoResponse{
					ID:           1,
					Name:         "name",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "spec",
					FarmID:       1,
					FarmInfo: pond.FarmInfo{
						ID:       1,
						Name:     "farm",
						Location: "loc",
						Owner:    "owner",
						Area:     "area",
					},
				}, nil)
			},
			want: want{
				body: `{"data":{"id":1,"name":"name","capacity":1,"depth":1,"water_quality":1,"species":"spec","farm":{"id":1,"name":"farm","location":"loc","owner":"owner","area":"area"}},"code":200,"message":"success"}`,
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
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetPondInfoByID(gomock.Any()).Return(pond.GetPondInfoResponse{}, nil).AnyTimes()
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
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
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
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetPondInfoByID(gomock.Any()).Return(pond.GetPondInfoResponse{}, fmt.Errorf("record not found"))
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
			mockFunc: func(pondDomain mock_pond.MockPondDomain) {
				pondDomain.EXPECT().GetPondInfoByID(gomock.Any()).Return(pond.GetPondInfoResponse{}, fmt.Errorf("some error"))
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
			pondDomain := mock_pond.NewMockPondDomain(mockCtrl)
			tt.mockFunc(*pondDomain)

			handler := PondHandler{
				domain:       pondDomain,
				timeoutInSec: tt.args.timeout,
			}

			r := httptest.NewRequest(http.MethodPost, "/pond/{id}", strings.NewReader(""))

			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)

			vars := map[string]string{
				"id": tt.body,
			}
			r = mux.SetURLVars(r, vars)

			w := httptest.NewRecorder()
			handler.GetByIDPondHandler(w, r)
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
