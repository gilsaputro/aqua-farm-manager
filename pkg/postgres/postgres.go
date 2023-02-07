package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PostgresConfig is list config to create postgres client
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// PostgresMethod is list all available method for postgres
type PostgresMethod interface {
	Insert(model *StatMetrics) error
	CheckStatExists(stat StatMetrics) bool
	UpdateStat(stat *StatMetrics) error
	GetStatRecodByKey(stat *StatMetrics) error
}

// Client is a wrapper for Postgres client
type Client struct {
	db *gorm.DB
}

// NewPostgresClient is func to create postgres client
func NewPostgresClient(config string) (PostgresMethod, error) {
	db, err := gorm.Open("postgres", config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Automatically create the table for the struct
	db.AutoMigrate(&Farms{}, &Ponds{}, &FarmPondsMapping{}, &StatMetrics{})
	return &Client{db: db}, nil
}

// CheckStatExists is func to check if the data is stat exist
func (c *Client) CheckStatExists(stat StatMetrics) bool {
	var count = int64(0)
	c.db.Model(stat).Where("key = ?", stat.Key).Count(&count).Limit(1)
	return count > 0
}

// Insert is func to insert data into table
func (c *Client) Insert(stat *StatMetrics) error {
	err := c.db.Create(stat).Error
	return err
}

// GetStatRecodByKey is func to get data into stat table using key
func (c *Client) GetStatRecodByKey(stat *StatMetrics) error {
	return c.db.Where("key = ?", stat.Key).First(stat).Error
}

// UpdateStat is func to update data into stat table
func (c *Client) UpdateStat(stat *StatMetrics) error {
	err := c.db.Model(stat).Where("key = ?", stat.Key).Update(StatMetrics{Request: stat.Request, UniqAgent: stat.UniqAgent}).Error
	return err
}
