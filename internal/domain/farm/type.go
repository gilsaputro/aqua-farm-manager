package farm

import (
	"errors"
)

// list Domain error
var (
	ErrDuplicateFarm = errors.New("Farm Already Exists")
	ErrInvalidFarm   = errors.New("Farm Is Not Exists")
	ErrExistsPonds   = errors.New("Cannot Delete Farm While Ponds Is Exists")
)

// CreateDomainRequest struct is list parameter for Create Farm domain
type CreateDomainRequest struct {
	Name     string
	Location string
	Owner    string
	Area     string
}

// CreateDomainResponse struct is list parameter response for Create Farm domain
type CreateDomainResponse struct {
	ID uint
}

// DeleteDomainRequest struct is list parameter for Delete Farm domain
type DeleteDomainRequest struct {
	Name string
	ID   uint
}

// DeleteDomainResponse struct is list parameter for Delete Farm domain
type DeleteDomainResponse struct {
	Name string
	ID   uint
}

// UpdateDomainRequest struct is list parameter for Update Farm domain
type UpdateDomainRequest struct {
	Name     string
	Location string
	Owner    string
	Area     string
}

// UpdateDomainResponse struct is list parameter response for Update Farm domain
type UpdateDomainResponse struct {
	ID       uint
	Name     string
	Location string
	Owner    string
	Area     string
}

// GetFarmInfoResponse struct is list parameter response for GetFarmInfoByID domain
type GetFarmInfoResponse struct {
	ID        uint
	Name      string
	Location  string
	Owner     string
	Area      string
	PondIDs   []uint
	PondInfos []PondInfo
}

// DeleteAllResponse struct is list parameter response for GetFarmInfoByID domain
type DeleteAllResponse struct {
	ID      uint
	Name    string
	PondIds []uint
}

// PondInfo struct is list parameter response for pond
type PondInfo struct {
	ID           uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
}
