package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"aqua-farm-manager/cmd/aqua-farm-manager/config"
	"aqua-farm-manager/internal/app"
	"aqua-farm-manager/internal/app/farm"
	"aqua-farm-manager/internal/app/middleware"
	"aqua-farm-manager/internal/app/stat"
	statdomain "aqua-farm-manager/internal/domain/stat"
	statinfra "aqua-farm-manager/internal/infrastructure/stat"
	"aqua-farm-manager/pkg/elasticsearch"
	"aqua-farm-manager/pkg/postgres"
	"aqua-farm-manager/pkg/redis"
	"aqua-farm-manager/pkg/vault"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Servcer is list configuration to run Server
type Server struct {
	cfg         config.Config
	vault       vault.VaultMethod
	redis       redis.RedisMethod
	postgres    postgres.PostgresMethod
	es          elasticsearch.ElasticSearchMethod
	middleware  middleware.Middleware
	statDomain  statdomain.StatDomain
	statInfra   statinfra.StatStore
	statHandler stat.StatHandler
	farmHandler farm.FarmHandler
}

// NewServer is func to create server with all configuration
func NewServer() {
	s := Server{}

	// ======== Init Dependencies Related ========
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

	// ======== Init Dependencies Infra ========
	{
		statinf := statinfra.NewStatStore(s.redis)
		s.statInfra = statinf
		log.Println("Init-NewStatStore")
	}

	// ======== Init Dependencies Domain ========
	{
		statDom := statdomain.NewStatDomain(s.statInfra)
		s.statDomain = statDom
		log.Println("Init-NewStatDomain")
	}

	// ======== Init Dependencies Handler/App ========
	// Init Middleware
	{
		mdl := middleware.NewMiddleware(s.statDomain)
		s.middleware = mdl
		log.Println("Init-NewMiddleware")
	}

	// Init FarmHandler
	{
		var opts []farm.Option
		opts = append(opts, farm.WithTimeoutOptions(s.cfg.FarmHandler.TimeoutInSec))
		handler := farm.NewFarmHandler(opts...)

		log.Println("Init-FarmHandler")
		s.farmHandler = *handler
	}

	// Init StatHandler
	{
		var opts []stat.Option
		opts = append(opts, stat.WithTimeoutOptions(s.cfg.StatHandler.TimeoutInSec))
		handler := stat.NewStatHandler(s.statDomain, opts...)

		log.Println("Init-StatHandler")
		s.statHandler = *handler
	}

	r := mux.NewRouter()
	// Init Farm Path
	port := ":8080"
	farmPath := app.Farms
	r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.CreateFarmHandler)).Methods("POST")
	r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.GetFarmHandler)).Methods("GET")
	r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.UpdateFarmHandler)).Methods("PUT")
	r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.DeleteFarmHandler)).Methods("DELETE")

	statPath := app.Stat
	r.HandleFunc(statPath.String(), s.statHandler.GetStatHandler).Methods("GET")

	log.Println("Running On", port)
	http.ListenAndServe(port, r)
}

// Run is func to create server and invoke Start()
func Run() int {
	return 0
}
