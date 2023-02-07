package postgres

import "github.com/jinzhu/gorm"

// Farms struct to store farm information
type Farms struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Location string
	Owner    string
	Area     string
	Status   int
}

// Keyword struct to store keyword information
type Ponds struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Capacity     float64
	Depth        float64
	WaterQuality float64
	Species      string
	Type         int
	Status       int
}

// FarmPondsMapping struct to store FarmPondsMapping information
type FarmPondsMapping struct {
	FarmID  uint
	PondsID uint
}

// StatMetrics struct to store StatMetrics information
type StatMetrics struct {
	gorm.Model
	Key       string
	Request   int
	UniqAgent int
}
