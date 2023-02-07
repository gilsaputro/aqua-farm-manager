package stat

import (
	"aqua-farm-manager/pkg/postgres"
	"aqua-farm-manager/pkg/postgres/mock_postgres"
	"aqua-farm-manager/pkg/redis"
	"aqua-farm-manager/pkg/redis/mock_redis"
	"fmt"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewStatStore(t *testing.T) {
	type args struct {
		redis redis.RedisMethod
		pg    postgres.PostgresMethod
	}
	tests := []struct {
		name string
		args args
		want StatStore
	}{
		{
			name: "success flow",
			args: args{
				redis: &redis.Client{},
				pg:    &postgres.Client{},
			},
			want: &Stat{
				redis: &redis.Client{},
				pg:    &postgres.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatStore(tt.args.redis, tt.args.pg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStat_IngestMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		urlID  string
		method string
		ua     string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod)
		wantErr  bool
	}{
		{
			name: "success flow new ua",
			args: args{
				urlID:  "1",
				method: "GET",
				ua:     "abcdef",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().SETNX("P:1:GET:abcdef").Return(true, nil)
				//somehow on windows OS if mock inside go func is not called so set it to AnyTimes
				r.EXPECT().HINCRBY(gomock.Any(), gomock.Any()).AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "got error on SetNX",
			args: args{
				urlID:  "1",
				method: "GET",
				ua:     "abcdef",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().SETNX("P:1:GET:abcdef").Return(false, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := mock_redis.NewMockRedisMethod(mockCtrl)
			pg := mock_postgres.NewMockPostgresMethod(mockCtrl)

			tt.mockFunc(redis, pg)
			s := NewStatStore(redis, pg)
			if err := s.IngestMetrics(tt.args.urlID, tt.args.method, tt.args.ua); (err != nil) != tt.wantErr {
				t.Errorf("Stat.IngestMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStat_GetMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		urlID  string
		method string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod)
		want     Metrics
		want1    string
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().HGETALL("P:1:GET").Return(
					map[string]string{CountUA: "1", CountRequested: "2"},
					nil,
				)
			},

			want: Metrics{
				Request:   "2",
				UniqAgent: "1",
			},
			wantErr: false,
		},
		{
			name: "success flow with zero data",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().HGETALL("P:1:GET").Return(
					map[string]string{},
					nil,
				)
			},
			want: Metrics{
				Request:   "0",
				UniqAgent: "0",
			},
			wantErr: false,
		},
		{
			name: "error flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().HGETALL("P:1:GET").Return(
					map[string]string{},
					fmt.Errorf("some error"),
				)
			},
			want: Metrics{
				Request:   "0",
				UniqAgent: "0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := mock_redis.NewMockRedisMethod(mockCtrl)
			pg := mock_postgres.NewMockPostgresMethod(mockCtrl)

			tt.mockFunc(redis, pg)
			s := NewStatStore(redis, pg)
			got, err := s.GetMetrics(tt.args.urlID, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stat.GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Stat.GetMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStat_BackupMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		urlID     string
		method    string
		request   int
		uniqagent int
	}
	tests := []struct {
		name     string
		mockFunc func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod)
		args     args
		wantErr  bool
	}{
		{
			name: "success flow insert",
			args: args{
				urlID:     "1",
				method:    "Get",
				request:   5,
				uniqagent: 2,
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().CheckStatExists(postgres.StatMetrics{
					Key:       "P:1:Get",
					Request:   5,
					UniqAgent: 2,
				}).Return(false)

				p.EXPECT().Insert(&postgres.StatMetrics{
					Key:       "P:1:Get",
					Request:   5,
					UniqAgent: 2,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success flow update",
			args: args{
				urlID:     "1",
				method:    "Get",
				request:   5,
				uniqagent: 2,
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().CheckStatExists(postgres.StatMetrics{
					Key:       "P:1:Get",
					Request:   5,
					UniqAgent: 2,
				}).Return(true)

				p.EXPECT().UpdateStat(&postgres.StatMetrics{
					Key:       "P:1:Get",
					Request:   5,
					UniqAgent: 2,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "got error flow",
			args: args{
				urlID:     "1",
				method:    "Get",
				request:   5,
				uniqagent: 2,
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().CheckStatExists(postgres.StatMetrics{
					Key:       "P:1:Get",
					Request:   5,
					UniqAgent: 2,
				}).Return(true)

				p.EXPECT().UpdateStat(&postgres.StatMetrics{
					Key:       "P:1:Get",
					Request:   5,
					UniqAgent: 2,
				}).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := mock_redis.NewMockRedisMethod(mockCtrl)
			pg := mock_postgres.NewMockPostgresMethod(mockCtrl)

			tt.mockFunc(redis, pg)
			s := NewStatStore(redis, pg)
			if err := s.BackupMetrics(tt.args.urlID, tt.args.method, tt.args.request, tt.args.uniqagent); (err != nil) != tt.wantErr {
				t.Errorf("Stat.BackupMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStat_MigrateMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		urlID     string
		method    string
		request   string
		uniqagent string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod)
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				urlID:     "1",
				method:    "GET",
				request:   "2",
				uniqagent: "1",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				//somehow on windows OS if mock inside go func is not called so set it to AnyTimes
				r.EXPECT().HSET(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := mock_redis.NewMockRedisMethod(mockCtrl)
			pg := mock_postgres.NewMockPostgresMethod(mockCtrl)

			tt.mockFunc(redis, pg)
			s := NewStatStore(redis, pg)
			if err := s.MigrateMetrics(tt.args.urlID, tt.args.method, tt.args.request, tt.args.uniqagent); (err != nil) != tt.wantErr {
				t.Errorf("Stat.MigrateMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStat_GetStatData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		urlID  string
		method string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod)
		want     Metrics
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				stat := &postgres.StatMetrics{
					Key: "P:1:GET",
				}
				p.EXPECT().GetStatRecodByKey(stat).DoAndReturn(
					func(stat *postgres.StatMetrics) error {
						stat.Request = 2
						stat.UniqAgent = 1
						return nil
					})
			},
			want:    Metrics{Request: "2", UniqAgent: "1"},
			wantErr: false,
		},
		{
			name: "error flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				stat := &postgres.StatMetrics{
					Key: "P:1:GET",
				}
				p.EXPECT().GetStatRecodByKey(stat).Return(fmt.Errorf("some error"))
			},
			want:    Metrics{Request: "0", UniqAgent: "0"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := mock_redis.NewMockRedisMethod(mockCtrl)
			pg := mock_postgres.NewMockPostgresMethod(mockCtrl)

			tt.mockFunc(redis, pg)
			s := NewStatStore(redis, pg)
			got, err := s.GetStatData(tt.args.urlID, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stat.GetStatData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Stat.GetStatData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
