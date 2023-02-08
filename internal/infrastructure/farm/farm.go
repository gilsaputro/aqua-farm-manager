package farm

import (
	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	"errors"

	"github.com/jinzhu/gorm"
)

// FarmStore is set of methods for interacting with a farm storage system
type FarmStore interface {
	VerifyByName(name string) (bool, error)
	VerifyByID(Id uint) (bool, error)
	Verify(r *FarmInfraInfo) (bool, error)
	GetFarmNameByID(id uint) (string, error)
	GetFarmIDByName(name string) (uint, error)
	Create(r *FarmInfraInfo) error
	Delete(r *FarmInfraInfo) error
	Update(r *FarmInfraInfo) error
	GetFarmByName(r *FarmInfraInfo) error
	GetFarmByID(r *FarmInfraInfo) error
	GetFarmWithPaging(r GetFarmWithPagingRequest) ([]FarmInfraInfo, error)
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
func (f *Farm) Create(r *FarmInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
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
		return err
	}

	r.ID = farm.Model.ID

	return err
}

// Create is func to store farm into database
func (f *Farm) Update(r *FarmInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Model: gorm.Model{
			ID: r.ID,
		},
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
		Status:   model.Active.Value(),
	}

	err = update(db, farm)
	if err != nil {
		return err
	}

	return err
}

// Delete is func to soft delete farm into database
func (f *Farm) Delete(r *FarmInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Model: gorm.Model{
			ID: r.ID,
		},
		Name: r.Name,
	}

	err = delete(db, farm)
	if err != nil {
		return err
	}

	return err
}

// GetFarmByName is func get farm info based on name in database
func (f *Farm) GetFarmByName(r *FarmInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Name: r.Name,
	}

	err = getFarmbyName(db, farm)

	r.Name = farm.Name
	r.Area = farm.Area
	r.ID = farm.Model.ID
	r.Location = farm.Location
	r.Owner = farm.Owner

	return err
}

// GetFarmByID is func get farm info based on id in database
func (f *Farm) GetFarmByID(r *FarmInfraInfo) error {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Model: gorm.Model{
			ID: r.ID,
		},
	}

	err = getFarmbyID(db, farm)

	r.Name = farm.Name
	r.Area = farm.Area
	r.ID = farm.Model.ID
	r.Location = farm.Location
	r.Owner = farm.Owner

	return err
}

// VerifyByName is func to check if farm already exists by name
func (f *Farm) VerifyByName(name string) (bool, error) {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return false, errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Name: name,
	}

	return checkFarmExistsByName(db, farm), err
}

// VerifyByName is func to check if farm already exists by id
func (f *Farm) VerifyByID(Id uint) (bool, error) {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return false, errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Model: gorm.Model{
			ID: Id,
		},
	}

	return checkFarmExistsByID(db, farm), err
}

// Verify is func to check if farm already exists based on id and name
func (f *Farm) Verify(r *FarmInfraInfo) (bool, error) {
	var err error
	db := f.pg.GetDB()
	if db == nil {
		return false, errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Model: gorm.Model{
			ID: r.ID,
		},
		Name: r.Name,
	}

	return checkFarmExists(db, farm), err
}

// GetFarmIDByName is func to get FarmID by name
func (f *Farm) GetFarmIDByName(name string) (uint, error) {
	var err error

	db := f.pg.GetDB()
	if db == nil {
		return 0, errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Name: name,
	}

	err = getFarmbyName(db, farm)

	return farm.Model.ID, err
}

// getFarmIDbyName func to get farm id by name
func getFarmbyName(db *gorm.DB, farm *postgres.Farms) error {
	return db.Where("name = ? AND Status = ?", farm.Name, model.Active.Value()).First(&farm).Error
}

// GetFarmIDByName is func to get FarmName by id
func (f *Farm) GetFarmNameByID(id uint) (string, error) {
	var err error

	db := f.pg.GetDB()
	if db == nil {
		return "", errors.New("Database Client is not init")
	}

	farm := &postgres.Farms{
		Model: gorm.Model{
			ID: id,
		},
	}

	err = getFarmbyID(db, farm)

	return farm.Name, err
}

// getFarmNamebyID func to get farm name by id
func getFarmbyID(db *gorm.DB, farm *postgres.Farms) error {
	return db.Where("id = ? and status = ?", farm.Model.ID, model.Active.Value()).First(&farm).Error
}

// insert is func to insert data farm into database
func insert(db *gorm.DB, data interface{}) error {
	return db.Create(data).Error
}

// update is func to update data farm in database
func update(db *gorm.DB, farm *postgres.Farms) error {
	return db.Model(farm).Where("name = ? AND id = ? and status = ?", farm.Name, farm.Model.ID, model.Active.Value()).Updates(farm).Error
}

// delete is func to soft delete data farm into database with update the status to inactive
func delete(db *gorm.DB, farm *postgres.Farms) error {
	return db.Model(farm).Where("name = ? AND id = ? and status = ?", farm.Name, farm.Model.ID, model.Active.Value()).Update("status", model.Inactive.Value()).Error
}

// checkFarmExistsByName is func to check if the data is exist by name
func checkFarmExistsByName(db *gorm.DB, farm *postgres.Farms) bool {
	var count = int64(0)
	db.Model(farm).Where("name = ? AND status = 1", farm.Name).Count(&count).Limit(1)
	return count > 0
}

// checkFarmExistsByID is func to check if the data is exist by id
func checkFarmExistsByID(db *gorm.DB, farm *postgres.Farms) bool {
	var count = int64(0)
	db.Model(farm).Where("id = ? AND status = 1", farm.Model.ID).Count(&count).Limit(1)
	return count > 0
}

// checkFarmExists is func to check if the data is exist by id and name
func checkFarmExists(db *gorm.DB, farm *postgres.Farms) bool {
	var count = int64(0)
	db.Model(farm).Where("id = ? AND  name = ? AND status = 1", farm.Model.ID, farm.Name).Count(&count).Limit(1)
	return count > 0
}

// GetFarmWithPaging is func to get all farm with paging
func (f *Farm) GetFarmWithPaging(r GetFarmWithPagingRequest) ([]FarmInfraInfo, error) {
	var list []FarmInfraInfo
	var err error

	db := f.pg.GetDB()
	if db == nil {
		return list, errors.New("Database Client is not init")
	}
	farms, err := getFarmsWithPaging(db, r.Cursor, r.Size)

	for _, farm := range farms {
		info := FarmInfraInfo{
			ID:       farm.ID,
			Name:     farm.Name,
			Location: farm.Location,
			Owner:    farm.Owner,
			Area:     farm.Area,
		}
		list = append(list, info)
	}

	return list, err
}

func getFarmsWithPaging(db *gorm.DB, cursor int, size int) ([]postgres.Farms, error) {
	var farms []postgres.Farms
	err := db.Where("status = ?", model.Active.Value()).Limit(size).Offset((cursor - 1) * size).Find(&farms).Error
	if err != nil {
		return nil, err
	}
	return farms, err
}
