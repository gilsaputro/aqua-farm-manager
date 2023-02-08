package stat

import (
	"aqua-farm-manager/internal/infrastructure/stat"
	"aqua-farm-manager/internal/infrastructure/stat/mock_stat"
	"fmt"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewStatDomain(t *testing.T) {
	type args struct {
		store stat.StatStore
	}
	tests := []struct {
		name string
		args args
		want StatDomain
	}{
		{
			name: "success flow",
			args: args{
				store: &stat.Stat{},
			},
			want: &Stat{store: &stat.Stat{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatDomain(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStat_GenerateStatAPI(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(r *mock_stat.MockStatStore)
		want     map[string]StatMetrics
	}{
		{
			name: "success flow",
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "GET",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "PUT",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "POST",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "DELETE",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "GET",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "2", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "PUT",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "2", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "POST",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "2", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "DELETE",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "2", NumSuccess: "1", NumError: "1"}, nil)
			},
			want: map[string]StatMetrics{
				"DELETE /farms": {1, 1, 1, 1},
				"DELETE /ponds": {2, 1, 1, 1},
				"GET /farms":    {1, 1, 1, 1},
				"GET /ponds":    {2, 1, 1, 1},
				"POST /farms":   {1, 1, 1, 1},
				"POST /ponds":   {2, 1, 1, 1},
				"PUT /farms":    {1, 1, 1, 1},
				"PUT /ponds":    {2, 1, 1, 1},
			},
		},
		{
			name: "partial error flow",
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "GET",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "PUT",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "POST",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "1",
					Method: "DELETE",
				}).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil)
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "GET",
				}).Return(stat.MetricsInfo{}, fmt.Errorf("some error"))
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "PUT",
				}).Return(stat.MetricsInfo{}, fmt.Errorf("some error"))
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "POST",
				}).Return(stat.MetricsInfo{}, fmt.Errorf("some error"))
				r.EXPECT().GetMetrics(stat.GetMetricsRequest{
					UrlID:  "2",
					Method: "DELETE",
				}).Return(stat.MetricsInfo{}, fmt.Errorf("some error"))
			},
			want: map[string]StatMetrics{
				"DELETE /farms": {1, 1, 1, 1},
				"GET /farms":    {1, 1, 1, 1},
				"POST /farms":   {1, 1, 1, 1},
				"PUT /farms":    {1, 1, 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			infra := mock_stat.NewMockStatStore(mockCtrl)

			tt.mockFunc(infra)
			s := NewStatDomain(infra)

			if got := s.GenerateStatAPI(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stat.GenerateStatAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestStat_IngestStatAPI(t *testing.T) {
// 	type args struct {
// 		path   string
// 		method string
// 		ua     string
// 		code   int
// 	}
// 	tests := []struct {
// 		name     string
// 		mockFunc func(r *mock_stat.MockStatStore)
// 		args     args
// 	}{
// 		{
// 			name: "success flow",
// 			args: args{
// 				path:   "/farms",
// 				method: "GET",
// 				ua:     "abc",
// 			},
// 			mockFunc: func(r *mock_stat.MockStatStore) {
// 				r.EXPECT().IngestMetrics("1", "GET", gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name: "got error flow",
// 			args: args{
// 				path:   "/farms",
// 				method: "GET",
// 				ua:     "abc",
// 			},
// 			mockFunc: func(r *mock_stat.MockStatStore) {
// 				r.EXPECT().IngestMetrics("1", "GET", gomock.Any()).Return(fmt.Errorf("some error"))
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockCtrl := gomock.NewController(t)
// 			defer mockCtrl.Finish()
// 			infra := mock_stat.NewMockStatStore(mockCtrl)

// 			tt.mockFunc(infra)
// 			s := NewStatDomain(infra)
// 			s.IngestStatAPI(tt.args.path, tt.args.method, tt.args.ua, tt.args.code)
// 		})
// 	}
// }

// func TestStat_BackUpStat(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		mockFunc func(r *mock_stat.MockStatStore)
// 	}{
// 		{
// 			name: "success flow",
// 			mockFunc: func(r *mock_stat.MockStatStore) {
// 				r.EXPECT().GetMetrics("1", gomock.Any()).Return(stat.Metrics{Request: "1", UniqAgent: "1"}, nil).Times(4)
// 				r.EXPECT().BackupMetrics("1", gomock.Any(), 1, 1).Return(nil).Times(4)
// 				r.EXPECT().GetMetrics("2", gomock.Any()).Return(stat.Metrics{Request: "1", UniqAgent: "2"}, nil).Times(4)
// 				r.EXPECT().BackupMetrics("2", gomock.Any(), 1, 2).Return(nil).Times(4)
// 			},
// 		},
// 		{
// 			name: "partial error flow",
// 			mockFunc: func(r *mock_stat.MockStatStore) {
// 				r.EXPECT().GetMetrics("1", gomock.Any()).Return(stat.Metrics{Request: "1", UniqAgent: "1"}, fmt.Errorf("some error")).Times(4)
// 				r.EXPECT().GetMetrics("2", gomock.Any()).Return(stat.Metrics{Request: "1", UniqAgent: "2"}, nil).Times(4)
// 				r.EXPECT().BackupMetrics("2", gomock.Any(), 1, 2).Return(fmt.Errorf("some error")).Times(4)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockCtrl := gomock.NewController(t)
// 			defer mockCtrl.Finish()
// 			infra := mock_stat.NewMockStatStore(mockCtrl)

// 			tt.mockFunc(infra)
// 			s := NewStatDomain(infra)
// 			s.BackUpStat()
// 		})
// 	}
// }

// func TestStat_MigrateStat(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		mockFunc func(r *mock_stat.MockStatStore)
// 	}{
// 		{
// 			name: "success flow",
// 			mockFunc: func(r *mock_stat.MockStatStore) {
// 				//somehow on windows OS if mock inside go func is not called so set it to AnyTimes
// 				r.EXPECT().GetStatData(gomock.Any(), gomock.Any()).Return(stat.Metrics{}, nil).AnyTimes()
// 				r.EXPECT().MigrateMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockCtrl := gomock.NewController(t)
// 			defer mockCtrl.Finish()
// 			infra := mock_stat.NewMockStatStore(mockCtrl)

// 			tt.mockFunc(infra)
// 			s := NewStatDomain(infra)
// 			s.MigrateStat()
// 		})
// 	}
// }
