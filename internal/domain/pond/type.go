package pond

import "errors"

// list Domain error
var (
	ErrDuplicatePond = errors.New("Pond Is Already Exists")
	ErrInvalidFarm   = errors.New("Farm Is Not Exists")
	ErrInvalidPond   = errors.New("Pond Is Not Exists")
)

// CreateDomainRequest struct is list parameter request for pond domain
type CreateDomainRequest struct {
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	FarmID       uint
}

// CreateDomainResponse struct is list parameter response for pond domain
type CreateDomainResponse struct {
	PondID uint
}

// PondInfo struct is list parameter response for pond domain
type PondInfo struct {
	ID           uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
}

// UpdateDomainRequest struct is list parameter for Update Farm domain
type UpdateDomainRequest struct {
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	FarmID       uint
}

// UpdateDomainResponse struct is list parameter response for Update Farm domain
type UpdateDomainResponse struct {
	ID           uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	FarmID       uint
}

// DeleteDomainRequest struct is list parameter for Delete Pond domain
type DeleteDomainRequest struct {
	Name string
	ID   uint
}

// DeleteDomainResponse struct is list parameter for Delete Pond domain
type DeleteDomainResponse struct {
	Name string
	ID   uint
}
