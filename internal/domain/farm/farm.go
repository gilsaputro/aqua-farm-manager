package farm

import (
	"aqua-farm-manager/internal/infrastructure/farm"
	"aqua-farm-manager/internal/infrastructure/pond"
)

// FarmDomain is list method for Farm domain
type FarmDomain interface {
	CreateFarmInfo(r CreateDomainRequest) (CreateDomainResponse, error)
	DeleteFarmInfo(r DeleteDomainRequest) (DeleteDomainResponse, error)
	UpdateFarmInfo(r UpdateDomainRequest) (UpdateDomainResponse, error)
	GetFarmInfoByID(ID uint) (GetFarmInfoResponse, error)
	GetFarm(size, cursor int) ([]GetFarmInfoResponse, int, error)
	DeleteFarmsWithDependencies(ID uint) (DeleteAllResponse, error)
}

// Stat is list dependencies stat domain
type Farm struct {
	pondstore pond.PondStore
	farmstore farm.FarmStore
}

// NewFarmDomain is func to generate FarmDomain interface
func NewFarmDomain(store farm.FarmStore, pondstore pond.PondStore) FarmDomain {
	return &Farm{
		farmstore: store,
		pondstore: pondstore,
	}
}

// CreateFarmInfo is func to store and validate request to create farm info
func (f *Farm) CreateFarmInfo(r CreateDomainRequest) (CreateDomainResponse, error) {
	var err error
	var res CreateDomainResponse

	exists, err := f.farmstore.Verify(&farm.FarmInfraInfo{
		Name: r.Name,
	})

	if err != nil {
		return res, err
	}

	if exists {
		return res, ErrDuplicateFarm
	}

	farmsInfra := mapCreateFarmInfoRequest(r)

	err = f.farmstore.Create(&farmsInfra)
	if err != nil {
		return res, err
	}

	res.ID = farmsInfra.ID

	return res, err
}

func mapCreateFarmInfoRequest(r CreateDomainRequest) farm.FarmInfraInfo {
	return farm.FarmInfraInfo{
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
	}
}

// DeleteFarmInfo is func to soft delete farm info in database
func (f *Farm) DeleteFarmInfo(r DeleteDomainRequest) (DeleteDomainResponse, error) {
	var err error
	var res DeleteDomainResponse
	var exists bool

	verify := farm.FarmInfraInfo{}
	if r.ID != 0 {
		verify.ID = r.ID
	} else if len(r.Name) != 0 {
		verify.Name = r.Name
	} else {
		return res, ErrInvalidFarm
	}

	exists, err = f.farmstore.Verify(&verify)

	if err != nil {
		return res, err
	}

	if !exists {
		return res, ErrInvalidFarm
	}

	ponds := f.farmstore.GetActivePondsInFarm(verify.ID)

	if len(ponds) > 0 {
		return res, ErrExistsPonds
	}

	err = f.farmstore.Delete(&farm.FarmInfraInfo{
		ID:   verify.ID,
		Name: verify.Name,
	})
	if err != nil {
		return res, err
	}

	res.ID = verify.ID
	res.Name = verify.Name

	return res, err
}

// UpdateFarmInfo is func to update farm info in database
func (f *Farm) UpdateFarmInfo(r UpdateDomainRequest) (UpdateDomainResponse, error) {
	var err error
	var res UpdateDomainResponse
	var exists bool

	exists, err = f.farmstore.Verify(&farm.FarmInfraInfo{
		Name: r.Name,
	})
	if err != nil {
		return res, err
	}

	farmsInfra := &farm.FarmInfraInfo{
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
	}
	if !exists {
		err = f.farmstore.Create(farmsInfra)
	} else {
		err = f.farmstore.GetFarmByName(farmsInfra)
		if err != nil {
			return res, err
		}

		// validate nil request
		if r.Location != "" {
			farmsInfra.Location = r.Location
		}
		if r.Owner != "" {
			farmsInfra.Owner = r.Owner
		}
		if r.Area != "" {
			farmsInfra.Area = r.Area
		}

		err = f.farmstore.Update(farmsInfra)
	}

	if err != nil {
		return res, err
	}

	return UpdateDomainResponse{
		ID:       farmsInfra.ID,
		Name:     farmsInfra.Name,
		Location: farmsInfra.Location,
		Owner:    farmsInfra.Owner,
		Area:     farmsInfra.Area,
	}, err
}

// GetFarmInfoByID is func to get farm info by id
func (f *Farm) GetFarmInfoByID(ID uint) (GetFarmInfoResponse, error) {
	var err error

	farm := &farm.FarmInfraInfo{
		ID: ID,
	}

	err = f.farmstore.GetFarmByID(farm)
	if err != nil {
		return GetFarmInfoResponse{}, err
	}

	ids, err := f.pondstore.GetPondIDbyFarmID(farm.ID)
	if err != nil {
		return GetFarmInfoResponse{}, err
	}
	var listPond []PondInfo
	for _, id := range ids {
		pond := &pond.PondInfraInfo{
			ID: id,
		}
		err := f.pondstore.GetPondByID(pond)
		if err != nil {
			continue
		}
		listPond = append(listPond, PondInfo{
			ID:           pond.ID,
			Name:         pond.Name,
			Capacity:     pond.Capacity,
			Depth:        pond.Depth,
			WaterQuality: pond.WaterQuality,
			Species:      pond.Species,
		})
	}
	return GetFarmInfoResponse{
		ID:        farm.ID,
		Name:      farm.Name,
		Location:  farm.Location,
		Owner:     farm.Owner,
		Area:      farm.Area,
		PondIDs:   ids,
		PondInfos: listPond,
	}, err
}

// GetFarm is func to get farm info by id
func (f *Farm) GetFarm(size, cursor int) ([]GetFarmInfoResponse, int, error) {
	var err error
	var list []GetFarmInfoResponse
	farmsInfra, err := f.farmstore.GetFarmWithPaging(
		farm.GetFarmWithPagingRequest{
			Size:   size,
			Cursor: cursor,
		})

	if err != nil {
		return list, 0, err
	}

	for _, farm := range farmsInfra {
		ids, err := f.pondstore.GetPondIDbyFarmID(farm.ID)
		if err != nil {
			return list, 0, err
		}
		info := GetFarmInfoResponse{
			ID:       farm.ID,
			Name:     farm.Name,
			Location: farm.Location,
			Owner:    farm.Owner,
			Area:     farm.Area,
			PondIDs:  ids,
		}

		list = append(list, info)
	}

	nextPage := cursor + 1
	if len(farmsInfra) < size {
		nextPage = 0
	}

	return list, nextPage, err
}

// DeleteFarmsWithDependencies is func to delete farms and all ponds dependencies
func (f *Farm) DeleteFarmsWithDependencies(ID uint) (DeleteAllResponse, error) {
	var err error
	var res DeleteAllResponse
	var exists bool

	verify := farm.FarmInfraInfo{ID: ID}

	exists, err = f.farmstore.Verify(&verify)

	if err != nil {
		return res, err
	}

	if !exists {
		return res, ErrInvalidFarm
	}

	ponds := f.farmstore.GetActivePondsInFarm(verify.ID)

	for _, p := range ponds {
		verifyPond := pond.PondInfraInfo{
			ID: p,
		}
		_, err = f.pondstore.Verify(&verifyPond)
		if err != nil {
			return res, err
		}
		err = f.pondstore.Delete(&pond.PondInfraInfo{
			ID:   verifyPond.ID,
			Name: verifyPond.Name,
		})
		if err != nil {
			return res, err
		}
	}

	err = f.farmstore.Delete(&farm.FarmInfraInfo{
		ID:   verify.ID,
		Name: verify.Name,
	})
	if err != nil {
		return res, err
	}

	res.ID = verify.ID
	res.Name = verify.Name
	res.PondIds = ponds

	return res, err
}
