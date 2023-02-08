package pond

// PondDomainRequest struct is list parameter request for pond domain
type PondDomainRequest struct {
	PondID       uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	Type         int
	FarmID       int
}

// PondResponse struct is list parameter response for pond domain
type PondResponse struct {
	PondID uint
}
