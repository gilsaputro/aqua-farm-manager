package farm

import (
	"aqua-farm-manager/internal/infrastructure/farm"
	"aqua-farm-manager/internal/infrastructure/farm/mock_farm"
	"aqua-farm-manager/internal/infrastructure/pond"
	"aqua-farm-manager/internal/infrastructure/pond/mock_pond"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewFarmDomain(t *testing.T) {
	type args struct {
		store     farm.FarmStore
		pondstore pond.PondStore
	}
	tests := []struct {
		name string
		args args
		want FarmDomain
	}{
		{
			name: "success flow",
			args: args{
				store:     &farm.Farm{},
				pondstore: &pond.Pond{},
			},
			want: &Farm{
				pondstore: &pond.Pond{},
				farmstore: &farm.Farm{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFarmDomain(tt.args.store, tt.args.pondstore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFarmDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_CreateFarmInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)
	type args struct {
		r CreateDomainRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     CreateDomainResponse
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().Create(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.ID = 1
						return nil
					})
			},
			args: args{
				r: CreateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			want: CreateDomainResponse{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "error when create",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				r: CreateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			want:    CreateDomainResponse{},
			wantErr: true,
		},
		{
			name: "farm exists",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
			},
			args: args{
				r: CreateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			want:    CreateDomainResponse{},
			wantErr: true,
		},
		{
			name: "error verify",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, fmt.Errorf("some error"))
			},
			args: args{
				r: CreateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			want:    CreateDomainResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmDomain(farmStore, pondStore)
			got, err := s.CreateFarmInfo(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.CreateFarmInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.CreateFarmInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_DeleteFarmInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)
	type args struct {
		r DeleteDomainRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     DeleteDomainResponse
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				r: DeleteDomainRequest{
					Name: "farm",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{})
				farmStore.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			want: DeleteDomainResponse{
				Name: "farm",
				ID:   1,
			},
			wantErr: false,
		},
		{
			name: "success flow by id",
			args: args{
				r: DeleteDomainRequest{
					ID: 1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{})
				farmStore.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			want: DeleteDomainResponse{
				Name: "farm",
				ID:   1,
			},
			wantErr: false,
		},
		{
			name: "error while delete flow",
			args: args{
				r: DeleteDomainRequest{
					Name: "farm",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{})
				farmStore.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
		{
			name: "error still have pond",
			args: args{
				r: DeleteDomainRequest{
					Name: "farm",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{1, 2, 3})
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
		{
			name: "farm not exists",
			args: args{
				r: DeleteDomainRequest{
					Name: "farm",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
		{
			name: "error verify",
			args: args{
				r: DeleteDomainRequest{
					Name: "farm",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, fmt.Errorf("some error"))
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmDomain(farmStore, pondStore)
			got, err := s.DeleteFarmInfo(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.DeleteFarmInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.DeleteFarmInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_UpdateFarmInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)
	type args struct {
		r UpdateDomainRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		want     UpdateDomainResponse
		wantErr  bool
	}{
		{
			name: "success create flow",
			args: args{
				r: UpdateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "Name"
						r.Area = "Area"
						r.Location = "Location"
						r.Owner = "Owner"
						return false, nil
					})
				farmStore.EXPECT().Create(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.ID = 1
						return nil
					})
			},
			want: UpdateDomainResponse{
				ID:       1,
				Name:     "Name",
				Location: "Location",
				Owner:    "Owner",
				Area:     "Area",
			},
			wantErr: false,
		},
		{
			name: "success update flow",
			args: args{
				r: UpdateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						return true, nil
					})
				farmStore.EXPECT().GetFarmByName(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.ID = 1
						r.Name = "Name"
						r.Area = "Area"
						r.Location = "Location"
						r.Owner = "Owner"
						return nil
					})
				farmStore.EXPECT().Update(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.ID = 1
						r.Name = "Name"
						r.Area = "Area"
						r.Location = "Location"
						r.Owner = "Owner"
						return nil
					})
			},
			want: UpdateDomainResponse{
				ID:       1,
				Name:     "Name",
				Location: "Location",
				Owner:    "Owner",
				Area:     "Area",
			},
			wantErr: false,
		},
		{
			name: "error update flow",
			args: args{
				r: UpdateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						return true, nil
					})
				farmStore.EXPECT().GetFarmByName(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.ID = 1
						r.Name = "Name"
						r.Area = "Area"
						r.Location = "Location"
						r.Owner = "Owner"
						return nil
					})
				farmStore.EXPECT().Update(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
		{
			name: "error get farm flow",
			args: args{
				r: UpdateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						return true, nil
					})
				farmStore.EXPECT().GetFarmByName(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
		{
			name: "error verify flow",
			args: args{
				r: UpdateDomainRequest{
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, fmt.Errorf("some error"))
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmDomain(farmStore, pondStore)
			got, err := s.UpdateFarmInfo(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.UpdateFarmInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.UpdateFarmInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_GetFarmInfoByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)
	type args struct {
		ID uint
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     GetFarmInfoResponse
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmByID(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.Name = "name"
						r.Area = "area"
						r.ID = 1
						r.Location = "location"
						r.Owner = "owner"
						return nil
					})
				pondStore.EXPECT().GetPondIDbyFarmID(gomock.Any()).Return([]uint{1}, nil)
				pondStore.EXPECT().GetPondByID(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					r.Name = "name"
					r.Capacity = 1
					r.Depth = 1
					r.WaterQuality = 1
					r.Species = "ikan"
					return nil
				})
			},
			args: args{
				ID: 1,
			},
			want: GetFarmInfoResponse{
				ID:       1,
				Name:     "name",
				Location: "location",
				Owner:    "owner",
				Area:     "area",
				PondIDs:  []uint{1},
				PondInfos: []PondInfo{
					{
						ID:           1,
						Name:         "name",
						Capacity:     1,
						Depth:        1,
						WaterQuality: 1,
						Species:      "ikan",
					},
				},
			},
		},
		{
			name: "error get pond id flow",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmByID(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.Name = "name"
						r.Area = "area"
						r.ID = 1
						r.Location = "location"
						r.Owner = "owner"
						return nil
					})
				pondStore.EXPECT().GetPondIDbyFarmID(gomock.Any()).Return([]uint{1}, fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    GetFarmInfoResponse{},
			wantErr: true,
		},
		{
			name: "error get farm flow",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmByID(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    GetFarmInfoResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmDomain(farmStore, pondStore)
			got, err := s.GetFarmInfoByID(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.GetFarmInfoByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.GetFarmInfoByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarm_GetFarm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)
	type args struct {
		size   int
		cursor int
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     []GetFarmInfoResponse
		want1    int
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmWithPaging(gomock.Any()).Return([]farm.FarmInfraInfo{
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
				}, nil)
				pondStore.EXPECT().GetPondIDbyFarmID(uint(1)).Return([]uint{1, 2}, nil)
				pondStore.EXPECT().GetPondIDbyFarmID(uint(2)).Return([]uint{3, 4}, nil)
			},
			args: args{
				size:   10,
				cursor: 1,
			},
			want: []GetFarmInfoResponse{
				{
					ID:       1,
					Name:     "1",
					Location: "1",
					Owner:    "1",
					Area:     "1",
					PondIDs:  []uint{1, 2},
				},
				{
					ID:       2,
					Name:     "2",
					Location: "2",
					Owner:    "2",
					Area:     "2",
					PondIDs:  []uint{3, 4},
				},
			},
			want1:   0,
			wantErr: false,
		},
		{
			name: "success flow with next cursor",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmWithPaging(gomock.Any()).Return([]farm.FarmInfraInfo{
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
				}, nil)
				pondStore.EXPECT().GetPondIDbyFarmID(uint(1)).Return([]uint{1, 2}, nil)
				pondStore.EXPECT().GetPondIDbyFarmID(uint(2)).Return([]uint{3, 4}, nil)
			},
			args: args{
				size:   2,
				cursor: 1,
			},
			want: []GetFarmInfoResponse{
				{
					ID:       1,
					Name:     "1",
					Location: "1",
					Owner:    "1",
					Area:     "1",
					PondIDs:  []uint{1, 2},
				},
				{
					ID:       2,
					Name:     "2",
					Location: "2",
					Owner:    "2",
					Area:     "2",
					PondIDs:  []uint{3, 4},
				},
			},
			want1:   2,
			wantErr: false,
		},
		{
			name: "error GetPondIDbyFarmID",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmWithPaging(gomock.Any()).Return([]farm.FarmInfraInfo{
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
				}, nil)
				pondStore.EXPECT().GetPondIDbyFarmID(gomock.Any()).Return([]uint{}, fmt.Errorf("some error"))
			},
			args: args{
				size:   10,
				cursor: 1,
			},
			want1:   0,
			wantErr: true,
		},
		{
			name: "Error GetFarmWithPaging",
			mockFunc: func() {
				farmStore.EXPECT().GetFarmWithPaging(gomock.Any()).Return([]farm.FarmInfraInfo{
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
				}, fmt.Errorf("some error"))
			},
			args: args{
				size:   10,
				cursor: 1,
			},
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmDomain(farmStore, pondStore)
			got, got1, err := s.GetFarm(tt.args.size, tt.args.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.GetFarm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.GetFarm() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Farm.GetFarm() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFarm_DeleteFarmsWithDependencies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)
	type args struct {
		ID uint
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     DeleteAllResponse
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{1})
				pondStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *pond.PondInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "pond"
						return true, nil
					})
				pondStore.EXPECT().Delete(gomock.Any()).Return(nil)
				farmStore.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			args: args{
				ID: 1,
			},
			want: DeleteAllResponse{
				ID:      1,
				Name:    "farm",
				PondIds: []uint{1},
			},
			wantErr: false,
		},
		{
			name: "error delete farm flow",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{1})
				pondStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *pond.PondInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "pond"
						return true, nil
					})
				pondStore.EXPECT().Delete(gomock.Any()).Return(nil)
				farmStore.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    DeleteAllResponse{},
			wantErr: true,
		},
		{
			name: "error pond flow",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{1})
				pondStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *pond.PondInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "pond"
						return true, nil
					})
				pondStore.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    DeleteAllResponse{},
			wantErr: true,
		},
		{
			name: "error verify pond flow",
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) (bool, error) {
						r.ID = 1
						r.Name = "farm"
						return true, nil
					})
				farmStore.EXPECT().GetActivePondsInFarm(uint(1)).Return([]uint{1})
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    DeleteAllResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewFarmDomain(farmStore, pondStore)
			got, err := s.DeleteFarmsWithDependencies(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Farm.DeleteFarmsWithDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Farm.DeleteFarmsWithDependencies() = %v, want %v", got, tt.want)
			}
		})
	}
}
