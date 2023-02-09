package farm

import (
	pondDomain "aqua-farm-manager/internal/domain/pond"
	"aqua-farm-manager/internal/infrastructure/farm"
	"aqua-farm-manager/internal/infrastructure/pond"
)

// FarmDomain is list method for Farm domain
type FarmDomain interface {
	CreateFarmInfo(r CreateDomainRequest) (CreateDomainResponse, error)
	DeleteFarmInfo(r DeleteDomainRequest) (DeleteDomainResponse, error)
	UpdateFarmInfo(r UpdateDomainRequest) (UpdateDomainResponse, error)
	GetFarmInfoByID(ID uint) (GetFarmInfoByIDResponse, error)
	GetFarm(size, cursor int) ([]GetFarmInfoByIDResponse, int, error)
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

	farm := mapCreateFarmInfoRequest(r)

	err = f.farmstore.Create(&farm)
	if err != nil {
		return res, err
	}

	res.ID = farm.ID

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
	if r.ID != 0 && len(r.Name) > 0 {
		verify.ID = r.ID
		verify.Name = r.Name
	} else if r.ID != 0 {
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

	req := &farm.FarmInfraInfo{
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
	}
	if !exists {
		err = f.farmstore.Create(req)
	} else {
		err = f.farmstore.GetFarmByName(req)
		if err != nil {
			return res, err
		}

		// validate nil request
		if r.Location != "" {
			req.Location = r.Location
		}
		if r.Owner != "" {
			req.Owner = r.Owner
		}
		if r.Area != "" {
			req.Area = r.Area
		}

		err = f.farmstore.Update(req)
	}

	if err != nil {
		return res, err
	}

	return UpdateDomainResponse{
		ID:       req.ID,
		Name:     req.Name,
		Location: req.Location,
		Owner:    req.Owner,
		Area:     req.Area,
	}, err
}

// GetFarmInfoByID is func to get farm info by id
func (f *Farm) GetFarmInfoByID(ID uint) (GetFarmInfoByIDResponse, error) {
	var err error

	farm := &farm.FarmInfraInfo{
		ID: ID,
	}

	err = f.farmstore.GetFarmByID(farm)
	if err != nil {
		return GetFarmInfoByIDResponse{}, err
	}

	ids, err := f.pondstore.GetPondIDbyFarmID(farm.ID)
	if err != nil {
		return GetFarmInfoByIDResponse{}, err
	}
	var listPond []pondDomain.PondInfo
	for _, id := range ids {
		pond := &pond.PondInfraInfo{
			ID: id,
		}
		err := f.pondstore.GetPondByID(pond)
		if err != nil {
			continue
		}
		listPond = append(listPond, pondDomain.PondInfo{
			ID:           pond.ID,
			Name:         pond.Name,
			Capacity:     pond.Capacity,
			Depth:        pond.Depth,
			WaterQuality: pond.WaterQuality,
			Species:      pond.Species,
		})
	}
	return GetFarmInfoByIDResponse{
		ID:        farm.ID,
		Name:      farm.Name,
		Location:  farm.Location,
		Owner:     farm.Owner,
		Area:      farm.Area,
		PondInfos: listPond,
	}, err
}

// GetFarm is func to get farm info by id
func (f *Farm) GetFarm(size, cursor int) ([]GetFarmInfoByIDResponse, int, error) {
	var err error
	var list []GetFarmInfoByIDResponse
	farms, err := f.farmstore.GetFarmWithPaging(
		farm.GetFarmWithPagingRequest{
			Size:   size,
			Cursor: cursor,
		})

	if err != nil {
		return list, 0, err
	}

	for _, farm := range farms {
		ids, err := f.pondstore.GetPondIDbyFarmID(farm.ID)
		if err != nil {
			continue
		}
		info := GetFarmInfoByIDResponse{
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
	if len(farms) < size {
		nextPage = 0
	}

	return list, nextPage, err
}
