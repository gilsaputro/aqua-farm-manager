package pond

import (
	"aqua-farm-manager/internal/infrastructure/farm"
	"aqua-farm-manager/internal/infrastructure/pond"
)

// PondDomain is list method for pond domain
type PondDomain interface {
	CreatePondInfo(r CreateDomainRequest) (CreateDomainResponse, error)
	UpdatePondInfo(r UpdateDomainRequest) (UpdateDomainResponse, error)
	DeletePondInfo(r DeleteDomainRequest) (DeleteDomainResponse, error)
	GetPondInfoByID(ID uint) (GetPondInfoResponse, error)
	GetAllPond(size, cursor int) ([]GetPondInfoResponse, int, error)
}

// Stat is list dependencies stat domain
type Pond struct {
	pondstore pond.PondStore
	farmstore farm.FarmStore
}

// NewPondDomain is func to generate PondDomain interface
func NewPondDomain(pondstore pond.PondStore, farmstore farm.FarmStore) PondDomain {
	return &Pond{
		pondstore: pondstore,
		farmstore: farmstore,
	}
}

// CreatePondInfo is func to update farm info in database
func (p *Pond) CreatePondInfo(r CreateDomainRequest) (CreateDomainResponse, error) {
	var err error
	var res CreateDomainResponse
	var exists bool
	exists, err = p.farmstore.Verify(
		&farm.FarmInfraInfo{
			ID: r.FarmID,
		})
	if err != nil {
		return res, err
	}

	if !exists {
		return res, ErrInvalidFarm
	}

	exists, err = p.pondstore.Verify(
		&pond.PondInfraInfo{
			Name: r.Name,
		})

	if err != nil {
		return res, err
	}

	if exists {
		return res, ErrDuplicatePond
	}

	pondinfra := mapPondRequest(r)

	ponds := p.farmstore.GetActivePondsInFarm(pondinfra.FarmID)
	if len(ponds) > 10 {
		return res, ErrMaxPond
	}

	err = p.pondstore.Create(pondinfra)
	if err != nil {
		return res, err
	}

	res.PondID = pondinfra.ID

	return res, err
}

func mapPondRequest(r CreateDomainRequest) *pond.PondInfraInfo {
	return &pond.PondInfraInfo{
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		FarmID:       r.FarmID,
	}
}

// UpdatePondInfo is func to update pond info in database
func (p *Pond) UpdatePondInfo(r UpdateDomainRequest) (UpdateDomainResponse, error) {
	var err error
	var res UpdateDomainResponse
	var existsPond bool

	verify := &pond.PondInfraInfo{
		Name: r.Name,
	}
	existsPond, err = p.pondstore.Verify(verify)

	if err != nil {
		return res, err
	}

	if r.FarmID > 0 {
		// verify farm id
		existsFarm, err := p.farmstore.Verify(
			&farm.FarmInfraInfo{
				ID: r.FarmID,
			})
		if err != nil {
			return res, err
		}
		// reject if the farmid is not exists
		if !existsFarm {
			return res, ErrInvalidFarm
		}
	}

	pondInfra := &pond.PondInfraInfo{
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		FarmID:       r.FarmID,
	}

	if !existsPond {
		if pondInfra.FarmID < 1 {
			return res, ErrInvalidFarm
		}
		ponds := p.farmstore.GetActivePondsInFarm(pondInfra.FarmID)
		if len(ponds) >= 10 {
			return res, ErrMaxPond
		}
		err = p.pondstore.Create(pondInfra)
	} else {
		pondInfra.ID = verify.ID
		err = p.pondstore.GetPondByID(pondInfra)
		if err != nil {
			return res, err
		}

		// validate nil request
		if r.Species != "" {
			pondInfra.Species = r.Species
		}
		if r.Capacity > 0 {
			pondInfra.Capacity = r.Capacity
		}
		if r.Depth > 0 {
			pondInfra.Depth = r.Depth
		}
		if r.WaterQuality > 0 {
			pondInfra.WaterQuality = r.WaterQuality
		}
		// check before modify farm
		if r.FarmID != pondInfra.FarmID && r.FarmID != 0 {
			ponds := p.farmstore.GetActivePondsInFarm(r.FarmID)
			if len(ponds) >= 10 {
				return res, ErrMaxPond
			}

			pondInfra.FarmID = r.FarmID
		}

		err = p.pondstore.Update(pondInfra)
	}

	if err != nil {
		return res, err
	}

	return UpdateDomainResponse{
		ID:           pondInfra.ID,
		Name:         pondInfra.Name,
		Capacity:     pondInfra.Capacity,
		Depth:        pondInfra.Depth,
		WaterQuality: pondInfra.WaterQuality,
		Species:      pondInfra.Species,
		FarmID:       pondInfra.FarmID,
	}, err
}

// DeletePondInfo is func to soft delete pond info in database
func (p *Pond) DeletePondInfo(r DeleteDomainRequest) (DeleteDomainResponse, error) {
	var err error
	var res DeleteDomainResponse
	var exists bool

	verify := pond.PondInfraInfo{}

	if r.ID != 0 {
		verify.ID = r.ID
	} else if len(r.Name) != 0 {
		verify.Name = r.Name
	} else {
		return res, ErrInvalidPond
	}

	exists, err = p.pondstore.Verify(&verify)

	if err != nil {
		return res, err
	}

	if !exists {
		return res, ErrInvalidPond
	}

	err = p.pondstore.Delete(&pond.PondInfraInfo{
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

// GetPondInfoByID is func to get farm info by id
func (p *Pond) GetPondInfoByID(ID uint) (GetPondInfoResponse, error) {
	var err error

	pondInfra := &pond.PondInfraInfo{
		ID: ID,
	}

	err = p.pondstore.GetPondByID(pondInfra)

	if err != nil {
		return GetPondInfoResponse{}, err
	}

	farmInfra := &farm.FarmInfraInfo{
		ID: pondInfra.FarmID,
	}

	err = p.farmstore.GetFarmByID(farmInfra)

	if err != nil {
		return GetPondInfoResponse{}, err
	}

	return GetPondInfoResponse{
		ID:           pondInfra.ID,
		Name:         pondInfra.Name,
		Capacity:     pondInfra.Capacity,
		Depth:        pondInfra.Depth,
		WaterQuality: pondInfra.WaterQuality,
		Species:      pondInfra.Species,
		FarmInfo: FarmInfo{
			ID:       farmInfra.ID,
			Name:     farmInfra.Name,
			Location: farmInfra.Location,
			Owner:    farmInfra.Owner,
			Area:     farmInfra.Area,
		},
	}, err
}

// GetAllPond is func to get farm info by id
func (p *Pond) GetAllPond(size, cursor int) ([]GetPondInfoResponse, int, error) {
	var err error
	var list []GetPondInfoResponse
	pondInfra, err := p.pondstore.GetPondWithPaging(
		pond.GetPondWithPagingRequest{
			Size:   size,
			Cursor: cursor,
		})

	if err != nil {
		return list, 0, err
	}

	for _, pond := range pondInfra {
		info := GetPondInfoResponse{
			ID:           pond.ID,
			Name:         pond.Name,
			Capacity:     pond.Capacity,
			Depth:        pond.Depth,
			WaterQuality: pond.WaterQuality,
			Species:      pond.Species,
			FarmID:       pond.FarmID,
		}

		list = append(list, info)
	}

	nextPage := cursor + 1
	if len(pondInfra) < size {
		nextPage = 0
	}

	return list, nextPage, err
}
