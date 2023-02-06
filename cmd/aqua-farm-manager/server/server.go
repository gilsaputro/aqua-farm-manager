package server

import (
	"fmt"
	"log"
	"os"

	"aqua-farm-manager/cmd/aqua-farm-manager/config"
	"aqua-farm-manager/internal/infrastructure/elasticsearch"
	"aqua-farm-manager/internal/infrastructure/postgres"
	"aqua-farm-manager/internal/infrastructure/redis"
	"aqua-farm-manager/pkg/vault"

	"github.com/joho/godotenv"
)

// Servcer is list configuration to run Server
type Server struct {
	cfg      config.Config
	vault    vault.VaultMethod
	redis    redis.RedisMethod
	postgres postgres.PostgresMethod
	es       elasticsearch.ElasticSearchMethod
}

// NewServer is func to create server with all configuration
func NewServer() {
	s := Server{}

	// Load Env File
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return
	}

	// Init Vault
	{
		token := os.Getenv("VAULT_TOKEN")
		if len(token) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_TOKEN")
			return
		}

		host := os.Getenv("VAULT_HOST")
		if len(host) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_HOST")
			return
		}

		vaultMethod, err := vault.NewVaultClient(token, host)
		if err != nil {
			fmt.Print("[Got Error]-Vault :", err)
		}
		s.vault = vaultMethod

		log.Println("Init-Vault")
	}

	// Get Config from yaml
	{
		secret, err := s.vault.GetConfig()
		if err != nil {
			fmt.Print("[Got Error]-Load Secret :", err)
		}
		cfg, err := config.GetConfig(secret)
		if err != nil {
			fmt.Print("[Got Error]-Load Config :", err)
		}
		s.cfg = cfg

		log.Println("LOAD-Config")
	}

	// Init RedisClient
	{
		redisMethod, err := redis.NewRedisClient(redis.RedisConfig{
			RedisHost:        s.cfg.Redis.RedisHost,
			Password:         s.cfg.Redis.RedisPassword,
			MaxIdleInSec:     s.cfg.Redis.MaxIdleInSec,
			IdleTimeoutInSec: s.cfg.Redis.IdleTimeoutInSec,
		})
		if err != nil {
			fmt.Print("[Got Error]-Redis :", err)
		}

		s.redis = redisMethod

		log.Println("Init-Redis")
	}

	// Init Postgres
	{
		postgresMethod, err := postgres.NewPostgresClient(s.cfg.Postgres.Config)
		if err != nil {
			fmt.Print("[Got Error]-Postgres :", err)
		}

		s.postgres = postgresMethod

		log.Println("Init-Postgres")
	}

	// Init ElasticSearch
	{
		esMethod, err := elasticsearch.CreateESClient(s.cfg.ES.Host)
		if err != nil {
			fmt.Print("[Got Error]-Postgres :", err)
		}

		s.es = esMethod
		log.Println("Init-ElasticSearch")
	}

}

// Run is func to create server and invoke Start()
func Run() int {
	return 0
}
