package stat

import (
	"aqua-farm-manager/internal/infrastructure/stat"
	"aqua-farm-manager/internal/infrastructure/stat/mock_stat"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
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
				"DELETE /v1/farms": {1, 1, 1, 1},
				"DELETE /v1/ponds": {2, 1, 1, 1},
				"GET /v1/farms":    {1, 1, 1, 1},
				"GET /v1/ponds":    {2, 1, 1, 1},
				"POST /v1/farms":   {1, 1, 1, 1},
				"POST /v1/ponds":   {2, 1, 1, 1},
				"PUT /v1/farms":    {1, 1, 1, 1},
				"PUT /v1/ponds":    {2, 1, 1, 1},
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
				r.EXPECT().GetStatData(gomock.Any()).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil).Times(4)
			},
			want: map[string]StatMetrics{
				"DELETE /v1/farms": {1, 1, 1, 1},
				"GET /v1/farms":    {1, 1, 1, 1},
				"POST /v1/farms":   {1, 1, 1, 1},
				"PUT /v1/farms":    {1, 1, 1, 1},
				"DELETE /v1/ponds": {1, 1, 1, 1},
				"GET /v1/ponds":    {1, 1, 1, 1},
				"POST /v1/ponds":   {1, 1, 1, 1},
				"PUT /v1/ponds":    {1, 1, 1, 1},
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

func TestStat_IngestStatAPI(t *testing.T) {
	type args struct {
		path   string
		method string
		ua     string
		code   int
	}
	tests := []struct {
		name     string
		mockFunc func(r *mock_stat.MockStatStore)
		args     args
	}{
		{
			name: "success flow",
			args: args{
				path:   "/v1/farms",
				method: "GET",
				ua:     "abc",
				code:   200,
			},
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().IngestMetrics(stat.IngestMetricsRequest{
					UrlID:     "1",
					Method:    "GET",
					UA:        "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532",
					IsSuccess: true,
				}).Return(nil)
			},
		},
		{
			name: "got error flow",
			args: args{
				path:   "/v1/farms",
				method: "GET",
				ua:     "abc",
			},
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().IngestMetrics(gomock.Any()).Return(fmt.Errorf("some error"))
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
			s.IngestStatAPI(IngestStatRequest{
				Path:   tt.args.path,
				Method: tt.args.method,
				Ua:     tt.args.ua,
				Code:   tt.args.code,
			})
		})
	}
}

func TestStat_BackUpStat(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(r *mock_stat.MockStatStore)
	}{
		{
			name: "success flow",
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().GetMetrics(gomock.Any()).Times(4)
				r.EXPECT().BackupMetrics(gomock.Any()).Return(nil).Times(4)
				r.EXPECT().GetMetrics(gomock.Any()).Return(stat.MetricsInfo{NumRequest: "1", NumUniqAgent: "1", NumSuccess: "1", NumError: "1"}, nil).Times(4)
				r.EXPECT().BackupMetrics(gomock.Any()).Return(nil).Times(4)
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
			s.BackUpStat()
		})
	}
}

func TestStat_MigrateStat(t *testing.T) {
	metric1 := stat.MetricsInfo{
		NumRequest:   "1",
		NumUniqAgent: "1",
		NumSuccess:   "1",
		NumError:     "1",
	}
	metric2 := stat.MetricsInfo{
		NumRequest:   "2",
		NumUniqAgent: "2",
		NumSuccess:   "2",
		NumError:     "2",
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	infra := mock_stat.NewMockStatStore(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func(r *mock_stat.MockStatStore)
	}{
		{
			name: "success flow",
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "1",
					Method: "GET",
				}).Return(metric1, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "1",
					Method: "POST",
				}).Return(metric1, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "1",
					Method: "PUT",
				}).Return(metric1, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "1",
					Method: "DELETE",
				}).Return(metric1, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "2",
					Method: "GET",
				}).Return(metric2, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "2",
					Method: "POST",
				}).Return(metric2, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "2",
					Method: "PUT",
				}).Return(metric2, nil)
				r.EXPECT().GetStatData(stat.GetStatDataRequest{
					UrlID:  "2",
					Method: "DELETE",
				}).Return(metric2, nil)

				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "1",
					Method:  "GET",
					Metrics: metric1,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "1",
					Method:  "POST",
					Metrics: metric1,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "1",
					Method:  "PUT",
					Metrics: metric1,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "1",
					Method:  "DELETE",
					Metrics: metric1,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "2",
					Method:  "GET",
					Metrics: metric2,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "2",
					Method:  "POST",
					Metrics: metric2,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "2",
					Method:  "PUT",
					Metrics: metric2,
				}).Return(nil)
				r.EXPECT().MigrateMetrics(stat.MigrateMetricsRequest{
					UrlID:   "2",
					Method:  "DELETE",
					Metrics: metric2,
				}).Return(nil)
			},
		},
		{
			name: "error migrate flow",
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().GetStatData(gomock.Any()).Return(metric1, nil).AnyTimes()

				r.EXPECT().MigrateMetrics(gomock.Any()).Return(fmt.Errorf("some error")).AnyTimes()
			},
		},
		{
			name: "error get stat flow",
			mockFunc: func(r *mock_stat.MockStatStore) {
				r.EXPECT().GetStatData(gomock.Any()).Return(metric1, fmt.Errorf("some error")).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(infra)
			s := NewStatDomain(infra)
			s.MigrateStat()
		})
	}
}
