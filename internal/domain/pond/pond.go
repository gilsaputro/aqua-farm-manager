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

	pondinfo := mapPondRequest(r)
	err = p.pondstore.Create(pondinfo)
	if err != nil {
		return res, err
	}

	res.PondID = pondinfo.ID

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
	var exists bool

	verify := &pond.PondInfraInfo{
		Name: r.Name,
	}
	exists, err = p.pondstore.Verify(verify)

	if err != nil {
		return res, err
	}

	// reject if the pond is not exists
	if !exists {
		return res, ErrInvalidPond
	}

	if r.FarmID > 0 {
		// verify farm id
		exists, err = p.farmstore.Verify(
			&farm.FarmInfraInfo{
				ID: r.FarmID,
			})
		if err != nil {
			return res, err
		}
		// reject if the farmid is not exists
		if !exists {
			return res, ErrInvalidFarm
		}
	}

	req := &pond.PondInfraInfo{
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		FarmID:       r.FarmID,
	}

	if !exists {
		err = p.pondstore.Create(req)
	} else {
		req.ID = verify.ID
		err = p.pondstore.GetPondByID(req)
		if err != nil {
			return res, err
		}

		// validate nil request
		if r.Species != "" {
			req.Species = r.Species
		}
		if r.Capacity > 0 {
			req.Capacity = r.Capacity
		}
		if r.Depth > 0 {
			req.Depth = r.Depth
		}
		if r.WaterQuality > 0 {
			req.WaterQuality = r.WaterQuality
		}
		if r.FarmID > 0 {
			req.FarmID = r.FarmID
		}

		err = p.pondstore.Update(req)
	}

	if err != nil {
		return res, err
	}

	return UpdateDomainResponse{
		ID:           req.ID,
		Name:         req.Name,
		Capacity:     req.Capacity,
		Depth:        req.Depth,
		WaterQuality: req.Depth,
		Species:      req.Species,
		FarmID:       req.FarmID,
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
		return res, ErrInvalidFarm
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
