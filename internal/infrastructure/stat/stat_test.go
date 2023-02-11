package stat

import (
	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	mock_postgres "aqua-farm-manager/pkg/postgres/mock"
	"aqua-farm-manager/pkg/redis"
	mock_redis "aqua-farm-manager/pkg/redis/mock"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
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
		urlID     string
		method    string
		ua        string
		isSuccess bool
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
				urlID:     "1",
				method:    "GET",
				ua:        "abcdef",
				isSuccess: true,
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().SETNX("P:1:GET:abcdef").Return(true, nil)
				r.EXPECT().HINCRBY("P:1:GET", CountUA).Return(nil)
				r.EXPECT().HINCRBY("P:1:GET", CountRequested).Return(nil)
				r.EXPECT().HINCRBY("P:1:GET", CountSuccess).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success flow is error metric",
			args: args{
				urlID:     "1",
				method:    "GET",
				ua:        "abcdef",
				isSuccess: false,
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				r.EXPECT().SETNX("P:1:GET:abcdef").Return(true, nil)
				r.EXPECT().HINCRBY("P:1:GET", CountUA).Return(nil)
				r.EXPECT().HINCRBY("P:1:GET", CountRequested).Return(nil)
				r.EXPECT().HINCRBY("P:1:GET", CountError).Return(nil)
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
			if err := s.IngestMetrics(
				IngestMetricsRequest{
					UrlID:     tt.args.urlID,
					Method:    tt.args.method,
					UA:        tt.args.ua,
					IsSuccess: tt.args.isSuccess,
				},
			); (err != nil) != tt.wantErr {
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
		want     MetricsInfo
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
					map[string]string{CountUA: "1", CountRequested: "2", CountError: "1", CountSuccess: "1"},
					nil,
				)
			},

			want: MetricsInfo{
				NumRequest:   "2",
				NumUniqAgent: "1",
				NumSuccess:   "1",
				NumError:     "1",
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
			want: MetricsInfo{
				NumRequest:   "0",
				NumUniqAgent: "0",
				NumSuccess:   "0",
				NumError:     "0",
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
			want: MetricsInfo{
				NumRequest:   "0",
				NumUniqAgent: "0",
				NumSuccess:   "0",
				NumError:     "0",
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
			got, err := s.GetMetrics(GetMetricsRequest{
				UrlID:  tt.args.urlID,
				Method: tt.args.method,
			})
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

var stat = &postgres.StatMetrics{
	Model: gorm.Model{
		ID: 123,
	},
	Key:        "P:1:GET",
	Request:    10,
	UniqAgent:  20,
	NumSuccess: 5,
	NumError:   2,
	Status:     model.Active.Value(),
}

var expectedRows = sqlmock.NewRows([]string{"id", "key", "request", "uniq_agent", "num_success", "num_error", "status"}).
	AddRow(stat.ID, stat.Key, stat.Request, stat.UniqAgent, stat.NumSuccess, stat.NumError, stat.Status)

func InitDBsMockupStat() (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)
	gormDB.LogMode(true)
	gormDB.SetLogger(log.New(os.Stdout, "\n", 0))
	gormDB.Debug()
	return db, mock, gormDB
}
func TestStat_GetStatData(t *testing.T) {
	db, mock, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
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
		want     MetricsInfo
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().GetDB().Return(gormDB)

				mock.ExpectQuery("SELECT (.+) FROM \"stat_metrics\"").
					WithArgs("P:1:GET", 1).
					WillReturnRows(expectedRows)
			},
			want: MetricsInfo{
				NumRequest:   "10",
				NumUniqAgent: "20",
				NumSuccess:   "5",
				NumError:     "2",
			},
			wantErr: false,
		},
		{
			name: "db not init flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().GetDB().Return(nil)
			},
			want: MetricsInfo{
				NumRequest:   "0",
				NumUniqAgent: "0",
				NumSuccess:   "0",
				NumError:     "0",
			},
			wantErr: true,
		},
		{
			name: "error flow",
			args: args{
				urlID:  "1",
				method: "GET",
			},
			mockFunc: func(r *mock_redis.MockRedisMethod, p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().GetDB().Return(gormDB)

				mock.ExpectQuery("SELECT (.+) FROM \"stat_metrics\"").
					WithArgs("P:1:GET", 1).WillReturnError(fmt.Errorf("some error"))
			},
			want: MetricsInfo{
				NumRequest:   "0",
				NumUniqAgent: "0",
				NumSuccess:   "0",
				NumError:     "0",
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
			got, err := s.GetStatData(GetStatDataRequest{
				UrlID:  tt.args.urlID,
				Method: tt.args.method,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Stat.GetStatData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Stat.GetStatData() got = %v, want %v", got, tt.want)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Stat.GetStatData() Expectation = %v", err)
			}
		})
	}
}

func TestStat_BackupMetrics(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		r BackupMetricsRequest
	}
	tests := []struct {
		name     string
		mockFunc func(p *mock_postgres.MockPostgresMethod)
		args     args
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func(p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().GetDB().Return(gormDB)

				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "stat_metrics" SET "status" = $1, "updated_at" = $2 WHERE "stat_metrics"."deleted_at" IS NULL AND ((key = $3 and status = $4))`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()

				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "stat_metrics" ("created_at","updated_at","deleted_at","key","request","uniq_agent","num_success","num_error","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "error nil db flow",
			mockFunc: func(p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().GetDB().Return(nil)
			},
			wantErr: true,
		},
		{
			name: "error on update flow",
			mockFunc: func(p *mock_postgres.MockPostgresMethod) {
				p.EXPECT().GetDB().Return(gormDB)

				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "stat_metrics" SET "status" = $1, "updated_at" = $2 WHERE "stat_metrics"."deleted_at" IS NULL AND ((key = $3 and status = $4))`)).WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := mock_postgres.NewMockPostgresMethod(mockCtrl)

			tt.mockFunc(pg)
			s := NewStatStore(&redis.Client{}, pg)
			if err := s.BackupMetrics(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Stat.BackupMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if the expected SQL statement was executed
			if err := mockDB.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestStat_MigrateMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		r MigrateMetricsRequest
	}
	tests := []struct {
		name     string
		mockFunc func(r *mock_redis.MockRedisMethod)
		args     args
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func(r *mock_redis.MockRedisMethod) {
				r.EXPECT().HSET(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(4)
			},
			args: args{
				r: MigrateMetricsRequest{
					UrlID:  "1",
					Method: "GET",
					Metrics: MetricsInfo{
						NumRequest:   "1",
						NumUniqAgent: "1",
						NumSuccess:   "1",
						NumError:     "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error flow",
			mockFunc: func(r *mock_redis.MockRedisMethod) {
				r.EXPECT().HSET(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error")).Times(4)
			},
			args: args{
				r: MigrateMetricsRequest{
					UrlID:  "1",
					Method: "GET",
					Metrics: MetricsInfo{
						NumRequest:   "1",
						NumUniqAgent: "1",
						NumSuccess:   "1",
						NumError:     "1",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := mock_redis.NewMockRedisMethod(mockCtrl)
			tt.mockFunc(redis)
			s := NewStatStore(redis, &postgres.Client{})
			if err := s.MigrateMetrics(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Stat.MigrateMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
