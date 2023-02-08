package farm

import "errors"

// list Domain error
var (
	ErrDuplicateFarm = errors.New("Farm Already Exists")
)

// FarmDomainRequest struct is list parameter for domain
type FarmDomainRequest struct {
	Name     string
	Location string
	Owner    string
	Area     string
}

// FarmDomainResponse struct is list parameter response for domain
type FarmDomainResponse struct {
	FarmID uint
}
