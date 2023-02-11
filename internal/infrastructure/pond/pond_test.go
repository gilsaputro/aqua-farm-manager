package pond

import (
	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	mock_postgres "aqua-farm-manager/pkg/postgres/mock"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewPondStore(t *testing.T) {
	type args struct {
		pg postgres.PostgresMethod
	}
	tests := []struct {
		name string
		args args
		want PondStore
	}{
		{
			name: "success",
			args: args{
				pg: &postgres.Client{},
			},
			want: &Pond{
				pg: &postgres.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPondStore(tt.args.pg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPondStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func InitDBsMockupStat() (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)
	gormDB.LogMode(true)
	gormDB.SetLogger(log.New(os.Stdout, "\n", 0))
	gormDB.Debug()
	return db, mock, gormDB
}

func TestPond_Create(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *PondInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ponds" ("created_at","updated_at","deleted_at","name","capacity","depth","water_quality","species","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()

				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "farm_ponds_mappings" ("created_at","updated_at","deleted_at","farm_id","ponds_id") VALUES ($1,$2,$3,$4,$5)`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
			r: &PondInfraInfo{
				ID:   1,
				Name: "a",
			},
			wantErr: false,
		},
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ponds" ("created_at","updated_at","deleted_at","name","capacity","depth","water_quality","species","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &PondInfraInfo{
				ID:   1,
				Name: "a",
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &PondInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "req nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			if err := s.Create(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Pond.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPond_GetPondIDbyFarmID(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		id       uint
		mockFunc func()
		wantErr  bool
		want     []uint
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farm_ponds_mappings" WHERE "farm_ponds_mappings"."deleted_at" IS NULL AND ((farm_id = $1))`)).WillReturnRows(sqlmock.NewRows([]string{"ponds_id"}).AddRow(1))

			},
			id:      1,
			want:    []uint{1},
			wantErr: false,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			id:      1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			got, err := s.GetPondIDbyFarmID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.GetPondIDbyFarmID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.GetPondIDbyFarmID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_GetPondByID(t *testing.T) {
	var pond1 = &postgres.Ponds{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "1",
		Capacity:     1,
		Depth:        1,
		WaterQuality: 1,
		Species:      "1",
		Status:       model.Active.Value(),
	}

	var expectedRows = sqlmock.NewRows([]string{"id", "name", "capacity", "depth", "water_quality", "species", "status"}).
		AddRow(pond1.ID, pond1.Name, pond1.Capacity, pond1.Depth, pond1.WaterQuality, pond1.Species, pond1.Status)
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	type args struct {
		r *PondInfraInfo
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ponds" WHERE "ponds"."deleted_at" IS NULL AND "ponds"."id" = $1 AND ((id = $2 and status = $3)) ORDER BY "ponds"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farm_ponds_mappings" WHERE "farm_ponds_mappings"."deleted_at" IS NULL AND ((ponds_id = $1)) ORDER BY "farm_ponds_mappings"."id" ASC LIMIT 1`)).WillReturnRows(sqlmock.NewRows([]string{"farms_id"}).AddRow(1))
			},
			args: args{
				r: &PondInfraInfo{
					ID: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ponds" WHERE "ponds"."deleted_at" IS NULL AND "ponds"."id" = $1 AND ((id = $2 and status = $3)) ORDER BY "ponds"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: args{
				r: &PondInfraInfo{
					ID: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "nil request",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				r: &PondInfraInfo{
					ID: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			if err := s.GetPondByID(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Pond.GetPondByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPond_GetPondByName(t *testing.T) {
	var pond1 = &postgres.Ponds{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "1",
		Capacity:     1,
		Depth:        1,
		WaterQuality: 1,
		Species:      "1",
		Status:       model.Active.Value(),
	}

	var expectedRows = sqlmock.NewRows([]string{"id", "name", "capacity", "depth", "water_quality", "species", "status"}).
		AddRow(pond1.ID, pond1.Name, pond1.Capacity, pond1.Depth, pond1.WaterQuality, pond1.Species, pond1.Status)
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	type args struct {
		r *PondInfraInfo
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ponds" WHERE "ponds"."deleted_at" IS NULL AND ((name = $1 and status = $2)) ORDER BY "ponds"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farm_ponds_mappings" WHERE "farm_ponds_mappings"."deleted_at" IS NULL AND ((ponds_id = $1)) ORDER BY "farm_ponds_mappings"."id" ASC LIMIT 1`)).WillReturnRows(sqlmock.NewRows([]string{"farms_id"}).AddRow(1))
			},
			args: args{
				r: &PondInfraInfo{
					Name: "1",
				},
			},
			wantErr: false,
		},
		{
			name: "error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ponds" WHERE "ponds"."deleted_at" IS NULL AND ((name = $1 and status = $2)) ORDER BY "ponds"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some errro"))
			},
			args: args{
				r: &PondInfraInfo{
					ID: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "nil request",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				r: &PondInfraInfo{
					Name: "1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			if err := s.GetPondByName(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Pond.GetPondByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPond_Update(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *PondInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "ponds" SET "capacity" = $1, "depth" = $2, "id" = $3, "name" = $4, "species" = $5, "status" = $6, "updated_at" = $7, "water_quality" = $8 WHERE "ponds"."deleted_at" IS NULL AND "ponds"."id" = $9 AND ((name = $10 AND id = $11 and status = $12))`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()

				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "farm_ponds_mappings" SET "farm_id" = $1, "updated_at" = $2 WHERE "farm_ponds_mappings"."deleted_at" IS NULL AND ((ponds_id = $3))`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			r: &PondInfraInfo{
				ID:           1,
				Name:         "1",
				Capacity:     1,
				Depth:        1,
				WaterQuality: 1,
				Species:      "1",
				FarmID:       1,
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "ponds" SET "capacity" = $1, "depth" = $2, "id" = $3, "name" = $4, "species" = $5, "status" = $6, "updated_at" = $7, "water_quality" = $8 WHERE "ponds"."deleted_at" IS NULL AND "ponds"."id" = $9 AND ((name = $10 AND id = $11 and status = $12))`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &PondInfraInfo{
				ID:           1,
				Name:         "1",
				Capacity:     1,
				Depth:        1,
				WaterQuality: 1,
				Species:      "1",
				FarmID:       1,
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &PondInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "req nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			if err := s.Update(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Pond.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPond_Delete(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *PondInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "ponds" SET "status" = $1, "updated_at" = $2 WHERE "ponds"."deleted_at" IS NULL AND "ponds"."id" = $3 AND ((name = $4 AND id = $5 and status = $6))`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			r: &PondInfraInfo{
				ID:   1,
				Name: "1",
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "ponds" SET "status" = $1, "updated_at" = $2 WHERE "ponds"."deleted_at" IS NULL AND "ponds"."id" = $3 AND ((name = $4 AND id = $5 and status = $6))`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &PondInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &PondInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "req nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			if err := s.Delete(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Pond.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPond_GetPondWithPaging(t *testing.T) {
	var pond1 = &postgres.Ponds{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "1",
		Capacity:     1,
		Depth:        1,
		WaterQuality: 1,
		Species:      "1",
		Status:       model.Active.Value(),
	}

	var pond2 = &postgres.Ponds{
		Model: gorm.Model{
			ID: 2,
		},
		Name:         "2",
		Capacity:     2,
		Depth:        2,
		WaterQuality: 2,
		Species:      "2",
		Status:       model.Active.Value(),
	}

	var expectedRows = sqlmock.NewRows([]string{"id", "name", "capacity", "depth", "water_quality", "species", "status", "farm_id"}).
		AddRow(pond1.ID, pond1.Name, pond1.Capacity, pond1.Depth, pond1.WaterQuality, pond1.Species, pond1.Status, 1).
		AddRow(pond2.ID, pond2.Name, pond2.Capacity, pond2.Depth, pond2.WaterQuality, pond2.Species, pond2.Status, 2)

	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	type args struct {
		r GetPondWithPagingRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		want     []PondInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				r: GetPondWithPagingRequest{
					Size:   2,
					Cursor: 1,
				},
			},
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT ponds.id, ponds.name, ponds.capacity, ponds.depth, ponds.water_quality, ponds.species, farm_ponds_mappings.farm_id FROM "ponds" left join farm_ponds_mappings on farm_ponds_mappings.ponds_id = ponds.id WHERE (status = $1) LIMIT 2 OFFSET 0`)).WillReturnRows(expectedRows)
			},

			want: []PondInfraInfo{
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
			},
		},
		{
			name: "db not init",
			args: args{
				r: GetPondWithPagingRequest{
					Size:   2,
					Cursor: 1,
				},
			},
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			got, err := s.GetPondWithPaging(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.GetPondWithPaging() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.GetPondWithPaging() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_Verify(t *testing.T) {
	var pond1 = &postgres.Ponds{
		Model: gorm.Model{
			ID: 1,
		},
		Name:         "1",
		Capacity:     1,
		Depth:        1,
		WaterQuality: 1,
		Species:      "1",
		Status:       model.Active.Value(),
	}

	var expectedRows = sqlmock.NewRows([]string{"id", "name", "capacity", "depth", "water_quality", "species", "status"}).
		AddRow(pond1.ID, pond1.Name, pond1.Capacity, pond1.Depth, pond1.WaterQuality, pond1.Species, pond1.Status)
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	type args struct {
		r *PondInfraInfo
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		wantErr  bool
		want     bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ponds"  WHERE "ponds"."deleted_at" IS NULL AND ((name = $1 AND Status = $2)) ORDER BY "ponds"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
			},
			args: args{
				r: &PondInfraInfo{
					Name: "1",
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "nil request",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				r: &PondInfraInfo{
					Name: "1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondStore(pg)
			got, err := s.Verify(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pond.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
