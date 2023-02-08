package farm

import (
	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	"errors"

	"github.com/jinzhu/gorm"
)

// FarmStore is set of methods for interacting with a farm storage system
type FarmStore interface {
	Verify(name string) (bool, error)
	Create(r FarmInfraRequest) (uint, error)
}

// Farm is list dependencies farm store
type Farm struct {
	pg postgres.PostgresMethod
}

// NewFarmStore is func to generate FarmStore interface
func NewFarmStore(pg postgres.PostgresMethod) FarmStore {
	return &Farm{
		pg: pg,
	}
}

// Create is func to store farm into database
func (f *Farm) Create(r FarmInfraRequest) (uint, error) {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return 0, errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
		Status:   model.Active.Value(),
	}

	err = insert(db, farm)
	if err != nil {
		return 0, err
	}

	return farm.ID, err
}

// Verify is func to check if farm already exists
func (f *Farm) Verify(name string) (bool, error) {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return false, errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Name: name,
	}

	return checkFarmExists(db, farm), err
}

// insert is func to insert data farm into database
func insert(db *gorm.DB, data interface{}) error {
	return db.Create(data).Error
}

// checkFarmExists is func to check if the data is exist
func checkFarmExists(db *gorm.DB, farm *postgres.Farms) bool {
	var count = int64(0)
	db.Model(farm).Where("name = ?", farm.Name).Count(&count).Limit(1)
	return count > 0
}
