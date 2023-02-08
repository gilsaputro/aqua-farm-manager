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
	GetDB() *gorm.DB
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

// GetDB is func to return database client
func (c *Client) GetDB() *gorm.DB {
	return c.db
}
