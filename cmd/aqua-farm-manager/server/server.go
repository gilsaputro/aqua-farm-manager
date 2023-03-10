package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"aqua-farm-manager/cmd/aqua-farm-manager/config"
	"aqua-farm-manager/internal/app"
	"aqua-farm-manager/internal/app/farm"
	"aqua-farm-manager/internal/app/middleware"
	"aqua-farm-manager/internal/app/pond"
	"aqua-farm-manager/internal/app/stat"
	"aqua-farm-manager/internal/app/trackingevent"
	farmdomain "aqua-farm-manager/internal/domain/farm"
	ponddomain "aqua-farm-manager/internal/domain/pond"
	statdomain "aqua-farm-manager/internal/domain/stat"
	farminfra "aqua-farm-manager/internal/infrastructure/farm"
	pondinfra "aqua-farm-manager/internal/infrastructure/pond"
	statinfra "aqua-farm-manager/internal/infrastructure/stat"
	"aqua-farm-manager/pkg/nsq"
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
	nsqProducer nsq.NsqMethod
	middleware  middleware.Middleware
	statDomain  statdomain.StatDomain
	statInfra   statinfra.StatStore
	statHandler stat.StatHandler
	farmDomain  farmdomain.FarmDomain
	farmInfra   farminfra.FarmStore
	farmHandler farm.FarmHandler
	pondDomain  ponddomain.PondDomain
	pondInfra   pondinfra.PondStore
	pondHandler pond.PondHandler
	httpServer  *http.Server
}

// NewServer is func to create server with all configuration
func NewServer() (*Server, error) {
	s := &Server{}

	// ======== Init Dependencies Related ========
	// Load Env File
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return s, err
	}

	// Init Vault
	{
		token := os.Getenv("VAULT_TOKEN")
		if len(token) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_TOKEN")
			return s, fmt.Errorf("[Got Error]-Vault Invalid VAULT_TOKEN")
		}

		host := os.Getenv("VAULT_HOST")
		if len(host) <= 0 {
			fmt.Print("[Got Error]-Vault Invalid VAULT_HOST")
			return s, fmt.Errorf("[Got Error]-Vault Invalid VAULT_HOST")
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

	// Init NSQ Producer
	{
		os.Setenv("NSQD_VERBOSE", "false")
		nsqProducer, err := nsq.NewNsqClient(s.cfg.NSQ.ProducerHost)
		if err != nil {
			fmt.Print("[Got Error]-NSQ Producer :", err)
		}
		s.nsqProducer = nsqProducer
		log.Println("Init-NSQ Producer")
	}

	// ======== Init Dependencies Infra ========
	// Init Stat Infra
	{
		statinf := statinfra.NewStatStore(s.redis, s.postgres)
		s.statInfra = statinf
		log.Println("Init-NewStatStore")
	}
	// Init Farm Infra
	{
		farmInf := farminfra.NewFarmStore(s.postgres)
		s.farmInfra = farmInf
		log.Println("Init-NewFarmStore")
	}
	// Init Pond Infra
	{
		pondInf := pondinfra.NewPondStore(s.postgres)
		s.pondInfra = pondInf
		log.Println("Init-NewPondStore")
	}

	// ======== Init Dependencies Domain ========
	// Init Stat Domain
	{
		statDom := statdomain.NewStatDomain(s.statInfra)
		s.statDomain = statDom
		log.Println("Init-NewStatDomain")
	}

	// Init Farm Domain
	{
		farmDom := farmdomain.NewFarmDomain(s.farmInfra, s.pondInfra)
		s.farmDomain = farmDom
		log.Println("Init-NewFarmDomain")
	}

	// Init Farm Domain
	{
		pondDom := ponddomain.NewPondDomain(s.pondInfra, s.farmInfra)
		s.pondDomain = pondDom
		log.Println("Init-NewPondDomain")
	}

	// ======== Init Dependencies Handler/App ========
	// Init Middleware
	{
		mdl := middleware.NewMiddleware(s.cfg.TrackingEvent.Topic, s.nsqProducer)
		s.middleware = mdl
		log.Println("Init-NewMiddleware")
	}

	// Init Stat Migrator
	{
		s.statDomain.MigrateStat()
		log.Println("Init-MigrateStatMetrics From Postgres To Redis")
	}

	// Init FarmHandler
	{
		var opts []farm.Option
		opts = append(opts, farm.WithTimeoutOptions(s.cfg.FarmHandler.TimeoutInSec))
		handler := farm.NewFarmHandler(s.farmDomain, opts...)

		log.Println("Init-FarmHandler")
		s.farmHandler = *handler
	}

	// Init PondHandler
	{
		var opts []pond.Option
		opts = append(opts, pond.WithTimeoutOptions(s.cfg.PondHandler.TimeoutInSec))
		handler := pond.NewPondHandler(s.pondDomain, opts...)

		log.Println("Init-FarmHandler")
		s.pondHandler = *handler
	}

	// Init StatHandler
	{
		var opts []stat.Option
		opts = append(opts, stat.WithTimeoutOptions(s.cfg.StatHandler.TimeoutInSec))
		handler := stat.NewStatHandler(s.statDomain, opts...)

		log.Println("Init-StatHandler")
		s.statHandler = *handler
	}

	// Init Tracking Event Consumer
	{
		consumer := trackingevent.NewTrackingEventConsumer(
			s.cfg.TrackingEvent.Topic,
			s.cfg.TrackingEvent.Channel,
			s.cfg.NSQ.ConsumerHost,
			s.cfg.TrackingEvent.MaxInFlight,
			s.cfg.TrackingEvent.NumConsumer,
			s.cfg.TrackingEvent.TimeoutInSec,
			s.statDomain)
		err := consumer.Start()
		if err != nil {
			fmt.Print("[Got Error]-NewTrackingEverntConsumer :", err)
		}
		log.Println("Init-TrackingEverntConsumer")
	}

	// Init Stat Backup Cron
	{
		s.statHandler.InitMigrate(s.cfg.StatHandler.BackupTimeInMinute)
		log.Println("Init-Backup Scheduler From Redis To Postgres")
	}

	// Init Router
	{
		r := mux.NewRouter()
		// Init Farm Path
		farmPath := app.Farms
		r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.CreateFarmHandler)).Methods("POST")
		r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.GetFarmHandler)).Methods("GET")
		r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.UpdateFarmHandler)).Methods("PUT")
		r.HandleFunc(farmPath.String(), s.middleware.Middleware(s.farmHandler.DeleteFarmHandler)).Methods("DELETE")

		// Init Farm Get By ID
		farmByIDPath := farmPath.String() + "/{id}"
		r.HandleFunc(farmByIDPath, s.middleware.Middleware(s.farmHandler.GetByIDFarmHandler)).Methods("GET")
		r.HandleFunc(farmByIDPath, s.middleware.Middleware(s.farmHandler.DeleteByIDFarmHandler)).Methods("DELETE")

		// Init Pond Path
		pondPath := app.Ponds
		r.HandleFunc(pondPath.String(), s.middleware.Middleware(s.pondHandler.CreatePondHandler)).Methods("POST")
		r.HandleFunc(pondPath.String(), s.middleware.Middleware(s.pondHandler.UpdatePondHandler)).Methods("PUT")
		r.HandleFunc(pondPath.String(), s.middleware.Middleware(s.pondHandler.DeletePondHandler)).Methods("Delete")
		r.HandleFunc(pondPath.String(), s.middleware.Middleware(s.pondHandler.GetPondHandler)).Methods("Get")

		// Init Pond Get By ID
		getPondByIDPath := pondPath.String() + "/{id}"
		r.HandleFunc(getPondByIDPath, s.middleware.Middleware(s.pondHandler.GetByIDPondHandler)).Methods("GET")

		// Init Stat Path
		statPath := app.Stat
		r.HandleFunc(statPath.String(), s.statHandler.GetStatHandler).Methods("GET")

		port := ":" + s.cfg.Port
		log.Println("running on port ", port)

		server := &http.Server{
			Addr:    port,
			Handler: r,
		}

		s.httpServer = server
	}
	return s, nil
}

func (s *Server) Start() int {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Wait for a signal to shut down the application
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Received interrupt signal, performing backup...")
	// Backup data from redis to postgres before shytdown
	s.statDomain.BackUpStat()
	// Create a context with a timeout to allow the server to cleanly shut down
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	log.Println("complete, shutting down.")
	return 0
}

// Run is func to create server and invoke Start()
func Run() int {
	s, err := NewServer()
	if err != nil {
		return 1
	}

	return s.Start()
}
