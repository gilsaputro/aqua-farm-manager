package pond

import "errors"

// list Domain error
var (
	ErrDuplicatePond = errors.New("Pond Is Already Exists")
	ErrInvalidFarm   = errors.New("Farm Is Not Exists")
)

// PondDomainRequest struct is list parameter request for pond domain
type PondDomainRequest struct {
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	FarmID       uint
}

// PondResponse struct is list parameter response for pond domain
type PondResponse struct {
	PondID uint
}
