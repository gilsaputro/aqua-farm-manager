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
	GetFarm(key string)
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
	db.AutoMigrate(&Farms{}, &Ponds{}, &FarmPondsMapping{})
	return &Client{db: db}, nil
}

func (c *Client) GetFarm(key string) {

}
