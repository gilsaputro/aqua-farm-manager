package pond

// PondInfraInfo struct is list parameter from Ponds Storage
type PondInfraInfo struct {
	ID           uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	FarmID       uint
}

// FarmPondsMapping is list parameter to store Ponds Farms Mapping Information
type FarmPondsMapping struct {
	FarmID  uint
	PondsID uint
}

//GetPondWithPagingRequest struct is list parameter to get all pond with page
type GetPondWithPagingRequest struct {
	Size   int
	Cursor int
}
