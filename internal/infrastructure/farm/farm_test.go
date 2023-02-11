package farm

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

func TestNewFarmStore(t *testing.T) {
	type args struct {
		pg postgres.PostgresMethod
	}
	tests := []struct {
		name string
		args args
		want FarmStore
	}{
		{
			name: "success",
			args: args{
				pg: &postgres.Client{},
			},
			want: &Farm{
				pg: &postgres.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFarmStore(tt.args.pg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFarmStore() = %v, want %v", got, tt.want)
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

func TestFarm_Create(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *FarmInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "farms" ("created_at","updated_at","deleted_at","name","location","owner","area","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "farms" ("created_at","updated_at","deleted_at","name","location","owner","area","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &FarmInfraInfo{
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
			s := NewFarmStore(pg)
			if err := s.Create(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Farm.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarm_Update(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *FarmInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "farms" SET "id" = $1, "status" = $2, "updated_at" = $3 WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $4 AND ((name = $5 AND id = $6 and status = $7))`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "farms" SET "id" = $1, "status" = $2, "updated_at" = $3 WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $4 AND ((name = $5 AND id = $6 and status = $7))`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &FarmInfraInfo{
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
			s := NewFarmStore(pg)
			if err := s.Update(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Farm.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarm_Delete(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *FarmInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "farms" SET "status" = $1, "updated_at" = $2 WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $3 AND ((name = $4 AND id = $5 and status = $6))`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "farms" SET "status" = $1, "updated_at" = $2 WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $3 AND ((name = $4 AND id = $5 and status = $6))`)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &FarmInfraInfo{
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
			s := NewFarmStore(pg)
			if err := s.Delete(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Farm.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarm_GetFarmByName(t *testing.T) {
	var farm1 = &postgres.Farms{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "Farm 1",
		Location: "California",
		Owner:    "John Doe",
		Area:     "100 acres",
		Status:   model.Active.Value(),
	}
	var oneRows = sqlmock.NewRows([]string{"id", "name", "location", "owner", "area", "status"}).
		AddRow(farm1.ID, farm1.Name, farm1.Location, farm1.Owner, farm1.Area, farm1.Status)
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *FarmInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farms" WHERE "farms"."deleted_at" IS NULL AND ((name = $1 AND Status = $2)) ORDER BY "farms"."id" ASC LIMIT 1`)).WillReturnRows(oneRows)
			},
			r: &FarmInfraInfo{
				Name: "Farm 1",
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farms" WHERE "farms"."deleted_at" IS NULL AND ((name = $1 AND Status = $2)) ORDER BY "farms"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &FarmInfraInfo{
				Name: "Farm 1",
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &FarmInfraInfo{
				Name: "Farm 1",
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
			s := NewFarmStore(pg)
			if err := s.GetFarmByName(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Farm.GetFarmByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarm_GetFarmByID(t *testing.T) {
	var farm1 = &postgres.Farms{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "Farm 1",
		Location: "California",
		Owner:    "John Doe",
		Area:     "100 acres",
		Status:   model.Active.Value(),
	}
	var oneRows = sqlmock.NewRows([]string{"id", "name", "location", "owner", "area", "status"}).
		AddRow(farm1.ID, farm1.Name, farm1.Location, farm1.Owner, farm1.Area, farm1.Status)
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *FarmInfraInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farms" WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $1 AND ((id = $2 AND Status = $3)) ORDER BY "farms"."id" ASC LIMIT 1`)).WillReturnRows(oneRows)
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "got error exec",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farms" WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $1 AND ((id = $2 AND Status = $3)) ORDER BY "farms"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &FarmInfraInfo{
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
			s := NewFarmStore(pg)
			if err := s.GetFarmByID(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Farm.GetFarmByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarm_Verify(t *testing.T) {
	var farm1 = &postgres.Farms{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "Farm 1",
		Location: "California",
		Owner:    "John Doe",
		Area:     "100 acres",
		Status:   model.Active.Value(),
	}
	var oneRows = sqlmock.NewRows([]string{"id", "name", "location", "owner", "area", "status"}).
		AddRow(farm1.ID, farm1.Name, farm1.Location, farm1.Owner, farm1.Area, farm1.Status)
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        *FarmInfraInfo
		wantErr  bool
		want     bool
	}{
		{
			name: "success with id",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farms"  WHERE "farms"."deleted_at" IS NULL AND "farms"."id" = $1 AND ((id = $2 AND Status = $3)) ORDER BY "farms"."id" ASC LIMIT 1`)).WillReturnRows(oneRows)
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: &FarmInfraInfo{
				ID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "req nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid req",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
			},
			r:       &FarmInfraInfo{},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmStore(pg)
			got, err := s.Verify(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Farm.Verify() got validation = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_GetFarmWithPaging(t *testing.T) {
	var farm1 = &postgres.Farms{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "1",
		Location: "1",
		Owner:    "1",
		Area:     "1",
		Status:   model.Active.Value(),
	}

	var farm2 = &postgres.Farms{
		Model: gorm.Model{
			ID: 2,
		},
		Name:     "2",
		Location: "2",
		Owner:    "2",
		Area:     "2",
		Status:   model.Active.Value(),
	}

	var expectedRows = sqlmock.NewRows([]string{"id", "name", "location", "owner", "area", "status"}).
		AddRow(farm1.ID, farm1.Name, farm1.Location, farm1.Owner, farm1.Area, farm1.Status).
		AddRow(farm2.ID, farm2.Name, farm2.Location, farm2.Owner, farm2.Area, farm2.Status)

	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        GetFarmWithPagingRequest
		wantErr  bool
		want     []FarmInfraInfo
	}{
		{
			name: "success with id",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "farms" WHERE "farms"."deleted_at" IS NULL AND ((status = $1)) LIMIT 2 OFFSET 0`)).WillReturnRows(expectedRows)
			},
			r: GetFarmWithPagingRequest{
				Size:   2,
				Cursor: 1,
			},
			wantErr: false,
			want: []FarmInfraInfo{
				{
					ID:       1,
					Name:     "1",
					Location: "1",
					Owner:    "1",
					Area:     "1",
				},
				{
					ID:       2,
					Name:     "2",
					Location: "2",
					Owner:    "2",
					Area:     "2",
				},
			},
		},
		{
			name: "db nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r: GetFarmWithPagingRequest{
				Size:   2,
				Cursor: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmStore(pg)
			got, err := s.GetFarmWithPaging(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.GetFarmWithPaging() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.GetFarmWithPaging() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_GetActivePondsInFarm(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		r        uint
		wantErr  bool
		want     []uint
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT "farm_ponds_mappings".* FROM "farm_ponds_mappings" JOIN ponds ON ponds.id = farm_ponds_mappings.ponds_id WHERE "farm_ponds_mappings"."deleted_at" IS NULL AND ((ponds.status = $1 AND farm_ponds_mappings.farm_id = $2))`)).WillReturnRows(sqlmock.NewRows([]string{"ponds_id"}).AddRow(1))
			},
			r:       1,
			wantErr: false,
			want:    []uint{1},
		},
		{
			name: "db nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			r:       1,
			wantErr: true,
			want:    []uint{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmStore(pg)
			if got := s.GetActivePondsInFarm(tt.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.GetActivePondsInFarm() = %v, want %v", got, tt.want)
			}
		})
	}
}
