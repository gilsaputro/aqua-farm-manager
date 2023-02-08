package pond

import (
	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	"errors"

	"github.com/jinzhu/gorm"
)

// StatStore is set of methods for interacting with a ponds storage system
type PondStore interface {
	Create(r PondRequest) (uint, error)
	VerifyByName(name string) (bool, error)
	GetPondIDbyFarmID(id uint) ([]uint, error)
	GetPondInfoByID(r *PondInfraInfo) error
}

// Pond is list dependencies pond store
type Pond struct {
	pg postgres.PostgresMethod
}

// NewPondStore is func to generate PondStore interface
func NewPondStore(pg postgres.PostgresMethod) PondStore {
	return &Pond{
		pg: pg,
	}
}

// Create is func to store ponds and mapping to database
func (p *Pond) Create(r PondRequest) (uint, error) {
	var err error
	db := p.pg.GetDB()
	if db == nil {
		return 0, errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		Status:       model.Active.Value(),
	}

	err = insert(db, pond)
	if err != nil {
		return 0, err
	}

	farmpondMapping := &postgres.FarmPondsMapping{
		FarmID:  r.FarmID,
		PondsID: pond.ID,
	}

	err = insert(db, farmpondMapping)
	if err != nil {
		return 0, err
	}

	return pond.ID, err
}

// VerifyByName is func to check if farm already exists by name
func (p *Pond) VerifyByName(name string) (bool, error) {
	var err error
	db := p.pg.GetDB()
	if db == nil {
		return false, errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{
		Name: name,
	}

	return checkFarmExistsByName(db, pond), err
}

// checkFarmExistsByName is func to check if the data is exist by name
func checkFarmExistsByName(db *gorm.DB, pond *postgres.Ponds) bool {
	var count = int64(0)
	db.Model(pond).Where("name = ? AND status = ?", pond.Name, model.Active.Value()).Count(&count).Limit(1)
	return count > 0
}

func insert(db *gorm.DB, data interface{}) error {
	return db.Create(data).Error
}

func (p *Pond) GetPondIDbyFarmID(id uint) ([]uint, error) {
	var list []uint
	var err error
	db := p.pg.GetDB()
	if db == nil {
		return list, errors.New("Database Client is not init")
	}

	mapping, err := getPondIDbyFarmID(db, id)

	for _, data := range mapping {
		list = append(list, data.PondsID)
	}

	return list, err
}

func getPondIDbyFarmID(db *gorm.DB, farmid uint) ([]postgres.FarmPondsMapping, error) {
	var mapping []postgres.FarmPondsMapping
	err := db.Where("farm_id = ?", farmid).Find(&mapping).Error
	return mapping, err
}

func (p *Pond) GetPondInfoByID(r *PondInfraInfo) error {
	var err error

	db := p.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{
		Model: gorm.Model{
			ID: r.ID,
		},
	}

	err = getPondByID(db, pond)

	r.ID = pond.Model.ID
	r.Name = pond.Name
	r.Capacity = pond.Capacity
	r.Depth = pond.Depth
	r.WaterQuality = pond.WaterQuality
	r.Species = pond.Species
	r.Status = pond.Status

	return err
}

func getPondByID(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Where("id = ? and status = ?", pond.Model.ID, model.Active.Value()).First(pond).Error
}
