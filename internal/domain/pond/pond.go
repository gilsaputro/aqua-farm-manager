package pond

import (
	"aqua-farm-manager/internal/infrastructure/farm"
	"aqua-farm-manager/internal/infrastructure/pond"
)

// PondDomain is list method for pond domain
type PondDomain interface {
	CreatePondInfo(r PondDomainRequest) (PondResponse, error)
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

func (p *Pond) CreatePondInfo(r PondDomainRequest) (PondResponse, error) {
	var err error
	var res PondResponse
	var exists bool
	exists, err = p.farmstore.VerifyByID(uint(r.FarmID))
	if err != nil {
		return res, err
	}

	if !exists {
		return res, ErrInvalidFarm
	}

	exists, err = p.pondstore.VerifyByName(r.Name)

	if exists {
		return res, ErrDuplicatePond
	}

	if err != nil {
		return res, err
	}

	if exists {
		return res, ErrDuplicatePond
	}

	infrarequest := mapPondRequest(r)
	pondID, err := p.pondstore.Create(infrarequest)
	if err != nil {
		return res, err
	}

	res.PondID = pondID

	return res, err
}

func mapPondRequest(r PondDomainRequest) pond.PondRequest {
	return pond.PondRequest{
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		FarmID:       r.FarmID,
	}
}
