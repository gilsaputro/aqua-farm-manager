package pond

// PondRequest struct is list parameter to store Ponds Information to Postgres
type PondRequest struct {
	PondID       uint
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	Type         int
	Status       int
	FarmID       int
}

// FarmPondsMapping is list parameter to store Ponds Farms Mapping Information
type FarmPondsMapping struct {
	FarmID  uint
	PondsID uint
}

// PondStatus denotes the pond status
type PondStatus int

// The following constant are the know pond status
const (
	PondStatusDeleted  PondStatus = -1
	PondStatusUnknown  PondStatus = 0
	PondStatusActive   PondStatus = 1
	PondStatusInactive PondStatus = 2
)

// Value convert Pondstatus into int
func (status PondStatus) Value() int {
	return int(status)
}
