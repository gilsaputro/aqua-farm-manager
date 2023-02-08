package farm

// FarmInfraInfo struct is list parameter info for farm
type FarmInfraInfo struct {
	ID       uint
	Name     string
	Location string
	Owner    string
	Area     string
}

//GetFarmWithPagingRequest struct is list parameter to get farm with page
type GetFarmWithPagingRequest struct {
	Size   int
	Cursor int
}
