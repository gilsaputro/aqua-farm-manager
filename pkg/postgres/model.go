package postgres

import "github.com/jinzhu/gorm"

// Farms struct to store farm information
type Farms struct {
	gorm.Model
	Name     string
	Location string
	Owner    string
	Area     string
	Status   int
}

// Ponds struct to store ponds information
type Ponds struct {
	gorm.Model
	Name         string
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	Type         int
	Status       int
}

// FarmPondsMapping struct to store FarmPondsMapping information
type FarmPondsMapping struct {
	gorm.Model
	FarmID  uint
	PondsID uint
}

// StatMetrics struct to store StatMetrics information
type StatMetrics struct {
	gorm.Model
	Key        string
	Request    int
	UniqAgent  int
	NumSuccess int
	NumError   int
}
