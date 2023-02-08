package pond

// PondRequest struct is list parameter to store Ponds Information to Postgres
type PondRequest struct {
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	Status       int
	FarmID       uint
}

// PondInfraInfo struct is list parameter from Ponds Storage
type PondInfraInfo struct {
	ID           uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	Status       int
	FarmID       uint
}

// FarmPondsMapping is list parameter to store Ponds Farms Mapping Information
type FarmPondsMapping struct {
	FarmID  uint
	PondsID uint
}
