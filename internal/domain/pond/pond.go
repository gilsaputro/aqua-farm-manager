package pond

import (
	"aqua-farm-manager/internal/infrastructure/pond"
)

// PondDomain is list method for pond domain
type PondDomain interface {
	CreatePondInfo(r PondDomainRequest) (PondResponse, error)
}

// Stat is list dependencies stat domain
type Pond struct {
	store pond.PondStore
}

// NewStatDomain is func to generat StatDomain interface
func NewStatDomain(store pond.PondStore) PondDomain {
	return &Pond{
		store: store,
	}
}

func (p *Pond) CreatePondInfo(r PondDomainRequest) (PondResponse, error) {
	var err error
	var res PondResponse

	infrarequest := mapPondRequest(r)

	pondID, err := p.store.Create(infrarequest)
	if err != nil {
		return res, err
	}

	res.PondID = pondID

	return res, err
}

func mapPondRequest(r PondDomainRequest) pond.PondRequest {
	return pond.PondRequest{
		PondID:       r.PondID,
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		Type:         r.Type,
		FarmID:       r.FarmID,
	}
}
