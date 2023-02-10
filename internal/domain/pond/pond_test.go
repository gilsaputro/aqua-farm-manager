package pond

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

func TestNewPondDomain(t *testing.T) {
	type args struct {
		pondstore pond.PondStore
		farmstore farm.FarmStore
	}
	tests := []struct {
		name string
		args args
		want PondDomain
	}{
		{
			name: "success",
			args: args{
				pondstore: &pond.Pond{},
				farmstore: &farm.Farm{},
			},
			want: &Pond{
				pondstore: &pond.Pond{},
				farmstore: &farm.Farm{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPondDomain(tt.args.pondstore, tt.args.farmstore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPondDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_CreatePondInfo(t *testing.T) {
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
			name: "success flow",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{})
				pondStore.EXPECT().Create(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					return nil
				})
			},
			want: CreateDomainResponse{
				PondID: 1,
			},
			wantErr: false,
		},
		{
			name: "error while create",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{})
				pondStore.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			want: CreateDomainResponse{
				PondID: 0,
			},
			wantErr: true,
		},
		{
			name: "error max pond",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
			},
			want: CreateDomainResponse{
				PondID: 0,
			},
			wantErr: true,
		},
		{
			name: "error duplicate pond",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
			},
			want: CreateDomainResponse{
				PondID: 0,
			},
			wantErr: true,
		},
		{
			name: "error verify pond",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, fmt.Errorf("some error"))
			},
			want: CreateDomainResponse{
				PondID: 0,
			},
			wantErr: true,
		},
		{
			name: "error farm not exists",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
			},
			want: CreateDomainResponse{
				PondID: 0,
			},
			wantErr: true,
		},
		{
			name: "error verify farm",
			args: args{
				r: CreateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "Ikan",
					FarmID:       1,
				},
			},
			mockFunc: func() {
				farmStore.EXPECT().Verify(gomock.Any()).Return(false, fmt.Errorf("some error"))
			},
			want: CreateDomainResponse{
				PondID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondDomain(pondStore, farmStore)
			got, err := s.CreatePondInfo(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.CreatePondInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.CreatePondInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_UpdatePondInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	farmStore := mock_farm.NewMockFarmStore(mockCtrl)
	pondStore := mock_pond.NewMockPondStore(mockCtrl)

	type args struct {
		r UpdateDomainRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     UpdateDomainResponse
		wantErr  bool
	}{
		{
			name: "success create flow",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{})
				pondStore.EXPECT().Create(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					r.Capacity = 1
					r.Depth = 1
					r.WaterQuality = 1
					r.Species = "ikan"
					r.FarmID = 1
					return nil
				})
			},
			args: args{
				r: UpdateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "ikan",
					FarmID:       1,
				},
			},
			want: UpdateDomainResponse{
				ID:           1,
				Name:         "Pond 1",
				Capacity:     1,
				Depth:        1,
				WaterQuality: 1,
				Species:      "ikan",
				FarmID:       1,
			},
			wantErr: false,
		},
		{
			name: "success update flow",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().GetPondByID(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					r.Capacity = 1
					r.Depth = 1
					r.WaterQuality = 1
					r.Species = "ikan"
					r.FarmID = 2
					return nil
				})
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{})
				pondStore.EXPECT().Update(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					r.Capacity = 1
					r.Depth = 1
					r.WaterQuality = 1
					r.Species = "ikan"
					r.FarmID = 1
					return nil
				})
			},
			args: args{
				r: UpdateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "ikan",
					FarmID:       1,
				},
			},
			want: UpdateDomainResponse{
				ID:           1,
				Name:         "Pond 1",
				Capacity:     1,
				Depth:        1,
				WaterQuality: 1,
				Species:      "ikan",
				FarmID:       1,
			},
			wantErr: false,
		},
		{
			name: "error max pondsuccess create flow",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			},
			args: args{
				r: UpdateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "ikan",
					FarmID:       1,
				},
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
		{
			name: "error max pond update flow",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().GetPondByID(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					r.Capacity = 1
					r.Depth = 1
					r.WaterQuality = 1
					r.Species = "ikan"
					r.FarmID = 2
					return nil
				})
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			},
			args: args{
				r: UpdateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "ikan",
					FarmID:       1,
				},
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
		{
			name: "error when create flow",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{})
				pondStore.EXPECT().Create(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					return fmt.Errorf("some error")
				})
			},
			args: args{
				r: UpdateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "ikan",
					FarmID:       1,
				},
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
		{
			name: "error when update flow",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				farmStore.EXPECT().Verify(gomock.Any()).Return(true, nil)
				pondStore.EXPECT().GetPondByID(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					r.ID = 1
					r.Capacity = 1
					r.Depth = 1
					r.WaterQuality = 1
					r.Species = "ikan"
					r.FarmID = 2
					return nil
				})
				farmStore.EXPECT().GetActivePondsInFarm(gomock.Any()).Return([]uint{})
				pondStore.EXPECT().Update(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) error {
					return fmt.Errorf("some error")
				})
			},
			args: args{
				r: UpdateDomainRequest{
					Name:         "Pond 1",
					Capacity:     1,
					Depth:        1,
					WaterQuality: 1,
					Species:      "ikan",
					FarmID:       1,
				},
			},
			want:    UpdateDomainResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondDomain(pondStore, farmStore)
			got, err := s.UpdatePondInfo(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.UpdatePondInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.UpdatePondInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_DeletePondInfo(t *testing.T) {
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
			name: "Success Flow By Name",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) (bool, error) {
					r.ID = 1
					r.Name = "P 1"
					return true, nil
				})
				pondStore.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			args: args{
				r: DeleteDomainRequest{
					ID: 1,
				},
			},
			want: DeleteDomainResponse{
				Name: "P 1",
				ID:   1,
			},
			wantErr: false,
		},
		{
			name: "Success Flow By ID",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) (bool, error) {
					r.ID = 1
					r.Name = "P 1"
					return true, nil
				})
				pondStore.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			args: args{
				r: DeleteDomainRequest{
					Name: "P 1",
				},
			},
			want: DeleteDomainResponse{
				Name: "P 1",
				ID:   1,
			},
			wantErr: false,
		},
		{
			name: "Pond Not Exists",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, nil)
			},
			args: args{
				r: DeleteDomainRequest{
					Name: "P 1",
				},
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
		{
			name: "error get pond info Exists",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).Return(false, fmt.Errorf("some error"))
			},
			args: args{
				r: DeleteDomainRequest{
					Name: "P 1",
				},
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
		{
			name: "Error Flow When Delete",
			mockFunc: func() {
				pondStore.EXPECT().Verify(gomock.Any()).DoAndReturn(func(r *pond.PondInfraInfo) (bool, error) {
					r.ID = 1
					r.Name = "P 1"
					return true, nil
				})
				pondStore.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				r: DeleteDomainRequest{
					Name: "P 1",
				},
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
		{
			name: "Invalid Req Delete",
			mockFunc: func() {

			},
			args: args{
				r: DeleteDomainRequest{},
			},
			want:    DeleteDomainResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondDomain(pondStore, farmStore)
			got, err := s.DeletePondInfo(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.DeletePondInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.DeletePondInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_GetPondInfoByID(t *testing.T) {
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
		want     GetPondInfoResponse
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				pondStore.EXPECT().GetPondByID(gomock.Any()).DoAndReturn(
					func(r *pond.PondInfraInfo) error {
						r.ID = 1
						r.Name = "P 1"
						r.Capacity = 1
						r.Depth = 1
						r.WaterQuality = 1
						r.Species = "ikan"
						r.FarmID = 1
						return nil
					})
				farmStore.EXPECT().GetFarmByID(gomock.Any()).DoAndReturn(
					func(r *farm.FarmInfraInfo) error {
						r.ID = 1
						r.Area = "Area"
						r.Location = "Location"
						r.Name = "Name"
						r.Owner = "Owner"
						return nil
					})
			},
			args: args{
				ID: 1,
			},
			want: GetPondInfoResponse{
				ID:           1,
				Name:         "P 1",
				Capacity:     1,
				Depth:        1,
				WaterQuality: 1,
				Species:      "ikan",
				FarmInfo: FarmInfo{
					ID:       1,
					Name:     "Name",
					Location: "Location",
					Owner:    "Owner",
					Area:     "Area",
				},
			},
			wantErr: false,
		},
		{
			name: "error when get farm info flow",
			mockFunc: func() {
				pondStore.EXPECT().GetPondByID(gomock.Any()).DoAndReturn(
					func(r *pond.PondInfraInfo) error {
						r.ID = 1
						r.Name = "P 1"
						r.Capacity = 1
						r.Depth = 1
						r.WaterQuality = 1
						r.Species = "ikan"
						r.FarmID = 1
						return nil
					})
				farmStore.EXPECT().GetFarmByID(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    GetPondInfoResponse{},
			wantErr: true,
		},
		{
			name: "error when get farm info flow",
			mockFunc: func() {
				pondStore.EXPECT().GetPondByID(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				ID: 1,
			},
			want:    GetPondInfoResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondDomain(pondStore, farmStore)
			got, err := s.GetPondInfoByID(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.GetPondInfoByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.GetPondInfoByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPond_GetAllPond(t *testing.T) {
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
		want     []GetPondInfoResponse
		want1    int
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				pondStore.EXPECT().GetPondWithPaging(gomock.Any()).Return(
					[]pond.PondInfraInfo{
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
					}, nil,
				)
			},
			args: args{
				size:   10,
				cursor: 1,
			},
			want: []GetPondInfoResponse{
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
			want1:   0,
			wantErr: false,
		},
		{
			name: "success flow",
			mockFunc: func() {
				pondStore.EXPECT().GetPondWithPaging(gomock.Any()).Return(
					[]pond.PondInfraInfo{
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
					}, nil,
				)
			},
			args: args{
				size:   2,
				cursor: 1,
			},
			want: []GetPondInfoResponse{
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
			want1:   2,
			wantErr: false,
		},
		{
			name: "success flow",
			mockFunc: func() {
				pondStore.EXPECT().GetPondWithPaging(gomock.Any()).Return(
					[]pond.PondInfraInfo{}, fmt.Errorf("some error"),
				)
			},
			args: args{
				size:   2,
				cursor: 1,
			},
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			s := NewPondDomain(pondStore, farmStore)
			got, got1, err := s.GetAllPond(tt.args.size, tt.args.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pond.GetAllPond() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pond.GetAllPond() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Pond.GetAllPond() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
