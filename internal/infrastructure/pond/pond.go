package pond

import (
	"aqua-farm-manager/pkg/postgres"
	"aqua-farm-manager/pkg/redis"
	"errors"

	"github.com/jinzhu/gorm"
)

// StatStore is set of methods for interacting with a ponds storage system
type PondStore interface {
	Create(r PondRequest) (uint, error)
}

// Pond is list dependencies pond store
type Pond struct {
	pg postgres.PostgresMethod
}

// NewPondStore is func to generate PondStore interface
func NewPondStore(redis redis.RedisMethod, pg postgres.PostgresMethod) PondStore {
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
		Type:         r.Type,
		Status:       PondStatusActive.Value(),
	}

	err = insert(db, pond)
	if err != nil {
		return 0, err
	}

	farmpondMapping := &postgres.FarmPondsMapping{
		FarmID:  uint(r.FarmID),
		PondsID: pond.ID,
	}

	err = insert(db, farmpondMapping)
	if err != nil {
		return 0, err
	}

	return pond.ID, err
}

func insert(db *gorm.DB, data interface{}) error {
	return db.Create(data).Error
}
