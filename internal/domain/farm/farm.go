package farm

import (
	"aqua-farm-manager/internal/infrastructure/farm"
)

// FarmDomain is list method for Farm domain
type FarmDomain interface {
	CreateFarmInfo(r FarmDomainRequest) (FarmDomainResponse, error)
}

// Stat is list dependencies stat domain
type Farm struct {
	store farm.FarmStore
}

// NewFarmDomain is func to generat FarmDomain interface
func NewFarmDomain(store farm.FarmStore) FarmDomain {
	return &Farm{
		store: store,
	}
}

func (p *Farm) CreateFarmInfo(r FarmDomainRequest) (FarmDomainResponse, error) {
	var err error
	var res FarmDomainResponse

	val, err := p.store.Verify(r.Name)

	if err != nil {
		return res, err
	}

	if val {
		return res, ErrDuplicateFarm
	}

	infrarequest := mapCreateFarmInfoRequest(r)

	FarmID, err := p.store.Create(infrarequest)
	if err != nil {
		return res, err
	}

	res.FarmID = FarmID

	return res, err
}

func mapCreateFarmInfoRequest(r FarmDomainRequest) farm.FarmInfraRequest {
	return farm.FarmInfraRequest{
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
	}
}
