package pond

import (
	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	"errors"

	"github.com/jinzhu/gorm"
)

// PondStore is set of methods for interacting with a ponds storage system
type PondStore interface {
	Verify(r *PondInfraInfo) (bool, error)
	GetPondIDbyFarmID(id uint) ([]uint, error)
	GetPondByID(r *PondInfraInfo) error
	GetPondByName(r *PondInfraInfo) error
	Create(r *PondInfraInfo) error
	Update(r *PondInfraInfo) error
	Delete(r *PondInfraInfo) error
	GetPondWithPaging(r GetPondWithPagingRequest) ([]PondInfraInfo, error)
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
func (p *Pond) Create(r *PondInfraInfo) error {
	var err error
	db := p.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
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
		return err
	}

	farmpondMapping := &postgres.FarmPondsMapping{
		FarmID:  r.FarmID,
		PondsID: pond.ID,
	}

	err = insert(db, farmpondMapping)
	if err != nil {
		return err
	}

	r.ID = pond.Model.ID

	return err
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

// GetPondByID is func to get pond info in database by pond id
func (p *Pond) GetPondByID(r *PondInfraInfo) error {
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

	if err != nil {
		return err
	}

	mapping := &postgres.FarmPondsMapping{
		PondsID: pond.ID,
	}
	err = getFarmIDbyPondID(db, mapping)

	r.ID = pond.Model.ID
	r.Name = pond.Name
	r.Capacity = pond.Capacity
	r.Depth = pond.Depth
	r.WaterQuality = pond.WaterQuality
	r.Species = pond.Species
	r.FarmID = mapping.FarmID
	return err
}

func getPondByID(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Where("id = ? and status = ?", pond.Model.ID, model.Active.Value()).First(pond).Error
}

// GetPondByName is func to get pond info in database by pond name
func (p *Pond) GetPondByName(r *PondInfraInfo) error {
	var err error

	db := p.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{
		Name: r.Name,
	}

	err = getPondByName(db, pond)

	if err != nil {
		return err
	}

	mapping := &postgres.FarmPondsMapping{
		PondsID: pond.ID,
	}
	err = getFarmIDbyPondID(db, mapping)

	r.ID = pond.Model.ID
	r.Name = pond.Name
	r.Capacity = pond.Capacity
	r.Depth = pond.Depth
	r.WaterQuality = pond.WaterQuality
	r.Species = pond.Species
	r.FarmID = mapping.FarmID
	return err
}

func getPondByName(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Where("name = ? and status = ?", pond.Name, model.Active.Value()).First(pond).Error
}

// Update is func to store pond into database
func (f *Pond) Update(r *PondInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{
		Model: gorm.Model{
			ID: r.ID,
		},
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		Status:       model.Active.Value(),
	}

	err = update(db, pond)
	if err != nil {
		return err
	}

	farmpondMapping := &postgres.FarmPondsMapping{
		FarmID:  r.FarmID,
		PondsID: pond.ID,
	}

	err = updateMapping(db, farmpondMapping)
	if err != nil {
		return err
	}

	return err
}

// update is func to update data pond in database
func update(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Model(pond).Where("name = ? AND id = ? and status = ?", pond.Name, pond.Model.ID, model.Active.Value()).Updates(pond).Error
}

// updateMapping is func to update mapping data pond farm in database
func updateMapping(db *gorm.DB, mapping *postgres.FarmPondsMapping) error {
	return db.Model(mapping).Where("ponds_id = ?", mapping.PondsID).Updates(postgres.FarmPondsMapping{FarmID: mapping.FarmID}).Error
}

// Delete is func to soft delete farm into database
func (f *Pond) Delete(r *PondInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{
		Model: gorm.Model{
			ID: r.ID,
		},
		Name: r.Name,
	}

	err = delete(db, pond)
	if err != nil {
		return err
	}

	return err
}

// delete is func to soft delete data pond into database with update the status to inactive
func delete(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Model(pond).Where("name = ? AND id = ? and status = ?", pond.Name, pond.Model.ID, model.Active.Value()).Update("status", model.Inactive.Value()).Error
}

// Verify is func to check if pond already exists based on id or name
func (p *Pond) Verify(r *PondInfraInfo) (bool, error) {
	var exists bool
	db := p.pg.GetDB()
	if db == nil {
		return false, errors.New("Database Client is not init")
	}

	pond := &postgres.Ponds{}

	if r.ID > 0 {
		pond.Model.ID = r.ID
		getPondbyID(db, pond)
	} else if len(r.Name) > 0 {
		pond.Name = r.Name
		getPondbyName(db, pond)
	} else {
		return exists, errors.New("Invalid Parameter Request")
	}

	if len(pond.Name) == 0 || pond.ID <= 0 {
		return exists, nil
	}

	exists = true
	r.Name = pond.Name
	r.ID = pond.ID

	return exists, nil
}

// getPondbyName func to get pond by name
func getPondbyName(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Where("name = ? AND Status = ?", pond.Name, model.Active.Value()).First(&pond).Error
}

// getPondbyID func to get pond by id
func getPondbyID(db *gorm.DB, pond *postgres.Ponds) error {
	return db.Where("id = ? AND Status = ?", pond.Model.ID, model.Active.Value()).First(&pond).Error
}

// GetPondWithPaging is func to get all pond with paging
func (p *Pond) GetPondWithPaging(r GetPondWithPagingRequest) ([]PondInfraInfo, error) {
	var list []PondInfraInfo
	var err error

	db := p.pg.GetDB()
	if db == nil {
		return list, errors.New("Database Client is not init")
	}
	ponds, err := getPondWithPaging(db, r.Cursor, r.Size)

	return ponds, err
}

func getPondWithPaging(db *gorm.DB, cursor int, size int) ([]PondInfraInfo, error) {
	var ponds []PondInfraInfo
	err := db.Table("ponds").
		Select("ponds.id, ponds.name, ponds.capacity, ponds.depth, ponds.water_quality, ponds.species, farm_ponds_mappings.farm_id").
		Joins("left join farm_ponds_mappings on farm_ponds_mappings.ponds_id = ponds.id").
		Where("status = ?", model.Active.Value()).
		Limit(size).
		Offset((cursor - 1) * size).
		Scan(&ponds).Error

	if err != nil {
		return nil, err
	}

	return ponds, nil
}

// getFarmIDbyPondID func to get farmid by pondid id
func getFarmIDbyPondID(db *gorm.DB, mapping *postgres.FarmPondsMapping) error {
	return db.Where("ponds_id = ?", mapping.PondsID).First(&mapping).Error
}
