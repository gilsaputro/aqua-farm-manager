package farm

import (
	"aqua-farm-manager/internal/infrastructure/farm"
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
	store farm.FarmStore
}

// NewFarmDomain is func to generate FarmDomain interface
func NewFarmDomain(store farm.FarmStore) FarmDomain {
	return &Farm{
		store: store,
	}
}

// CreateFarmInfo is func to store and validate request to create farm info
func (f *Farm) CreateFarmInfo(r CreateDomainRequest) (CreateDomainResponse, error) {
	var err error
	var res CreateDomainResponse

	val, err := f.store.VerifyByName(r.Name)

	if err != nil {
		return res, err
	}

	if val {
		return res, ErrDuplicateFarm
	}

	farm := mapCreateFarmInfoRequest(r)

	err = f.store.Create(&farm)
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
	if r.ID != 0 && len(r.Name) > 0 {
		exists, err = f.store.Verify(&farm.FarmInfraInfo{
			ID:   r.ID,
			Name: r.Name,
		})
	} else if r.ID != 0 {
		var name string
		exists, err = f.store.VerifyByID(r.ID)
		if err != nil {
			return res, err
		}
		name, err = f.store.GetFarmNameByID(r.ID)
		r.Name = name
	} else if len(r.Name) != 0 {
		var id uint
		exists, err = f.store.VerifyByName(r.Name)
		if err != nil {
			return res, err
		}
		id, err = f.store.GetFarmIDByName(r.Name)
		r.ID = id
	} else {
		return res, ErrInvalidFarm
	}

	if err != nil {
		return res, err
	}

	if !exists {
		return res, ErrInvalidFarm
	}
	res.ID = r.ID
	res.Name = r.Name
	err = f.store.Delete(&farm.FarmInfraInfo{
		ID:   r.ID,
		Name: r.Name,
	})
	if err != nil {
		return res, err
	}

	return res, err
}

// UpdateFarmInfo is func to update farm info in database
func (f *Farm) UpdateFarmInfo(r UpdateDomainRequest) (UpdateDomainResponse, error) {
	var err error
	var res UpdateDomainResponse
	var exists bool

	exists, err = f.store.VerifyByName(r.Name)
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
		err = f.store.Create(req)
	} else {
		err = f.store.GetFarmByName(req)
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

		err = f.store.Update(req)
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

	err = f.store.GetFarmByID(farm)
	if err != nil {
		return GetFarmInfoByIDResponse{}, err
	}

	return GetFarmInfoByIDResponse{
		ID:       farm.ID,
		Name:     farm.Name,
		Location: farm.Location,
		Owner:    farm.Owner,
		Area:     farm.Area,
	}, err
}

// GetFarm is func to get farm info by id
func (f *Farm) GetFarm(size, cursor int) ([]GetFarmInfoByIDResponse, int, error) {
	var err error
	var list []GetFarmInfoByIDResponse
	farms, err := f.store.GetFarmWithPaging(
		farm.GetFarmWithPagingRequest{
			Size:   size,
			Cursor: cursor,
		})

	if err != nil {
		return list, 0, err
	}

	for _, farm := range farms {
		info := GetFarmInfoByIDResponse{
			ID:       farm.ID,
			Name:     farm.Name,
			Location: farm.Location,
			Owner:    farm.Owner,
			Area:     farm.Area,
		}

		list = append(list, info)
	}

	nextPage := cursor + 1
	if len(farms) < size {
		nextPage = 0
	}

	return list, nextPage, err
}
