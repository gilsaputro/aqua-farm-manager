package stat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"aqua-farm-manager/pkg/postgres"
	"aqua-farm-manager/pkg/redis"
)

var (
	PathKeyMetrics = "P:<urlID>:<method>"
	UAKeyMetrics   = "P:<urlID>:<method>:<ua>"

	CountUA        = "Count_UA"
	CountRequested = "Count_Req"
)

// StatStore is list method to redis
type StatStore interface {
	IngestMetrics(urlID, method, ua string) error
	GetMetrics(urlID, method string) (Metrics, error)
	BackupMetrics(urlID, method string, request, uniqagent int) error
	MigrateMetrics(url, method, request, uniqagent string) error
	GetStatData(urlID, method string) (Metrics, error)
}

type Metrics struct {
	Request   string
	UniqAgent string
}

// Stat is list dependencies stat store
type Stat struct {
	redis redis.RedisMethod
	pg    postgres.PostgresMethod
}

// NewStatStore is func to generate StatStore interface
func NewStatStore(redis redis.RedisMethod, pg postgres.PostgresMethod) StatStore {
	return &Stat{
		redis: redis,
		pg:    pg,
	}
}

// IngestMetrics is func to ingest api metrics to redis
func (s *Stat) IngestMetrics(urlID, method, ua string) error {
	var err error
	pathKey := generatePathKeyMetrics(urlID, method)
	uakey := generateUAKeyMetrics(urlID, method, ua)

	// set key to redis, if the key is exists it's not uniq
	isNew, err := s.redis.SETNX(uakey)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	// incr count uniq ua if is new
	if isNew {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.redis.HINCRBY(pathKey, CountUA)
		}()
	}

	// incr count requested
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.redis.HINCRBY(pathKey, CountRequested)
	}()
	return err
}

// GetMetrics is func to get api metrics from redis
func (s *Stat) GetMetrics(urlID, method string) (Metrics, error) {
	var err error
	pathKey := generatePathKeyMetrics(urlID, method)

	metrics, err := s.redis.HGETALL(pathKey)
	if err != nil {
		return Metrics{Request: "0", UniqAgent: "0"}, err
	}

	cUA := metrics[CountUA]
	cReq := metrics[CountRequested]

	if len(cUA) == 0 {
		cUA = "0"
	}

	if len(cReq) == 0 {
		cReq = "0"
	}

	return Metrics{Request: cReq, UniqAgent: cUA}, nil
}

// BackupMetrics is func to backup metrics from redis to postgres
func (s *Stat) BackupMetrics(urlID, method string, request, uniqagent int) error {
	var err error
	pathKey := generatePathKeyMetrics(urlID, method)
	stat := postgres.StatMetrics{
		Key:       pathKey,
		Request:   request,
		UniqAgent: uniqagent,
	}

	val := s.pg.CheckStatExists(stat)
	if val {
		err = s.pg.UpdateStat(&stat)
	} else {
		err = s.pg.Insert(&stat)
	}

	return err
}

// MigrateMetrics is func to migrate metrics from postgres to redis
func (s *Stat) MigrateMetrics(urlID, method, request, uniqagent string) error {
	var errUA, errReq error
	key := generatePathKeyMetrics(urlID, method)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		errUA = s.redis.HSET(key, CountUA, uniqagent)
		if errUA != nil {
			log.Println("MigrateMetrics-Error Ingest Uniq Agent :", errUA)
		}
	}()

	go func() {
		defer wg.Done()
		errReq = s.redis.HSET(key, CountRequested, request)
		if errUA != nil {
			log.Println("MigrateMetrics-Error Ingest Requested :", errReq)
		}
	}()

	if errUA != nil || errReq != nil {
		return fmt.Errorf("got error while migrate")
	}

	return nil
}

// GetStatData is func to metrics from postgres
func (s *Stat) GetStatData(urlID, method string) (Metrics, error) {
	var err error
	pathKey := generatePathKeyMetrics(urlID, method)

	statMetrics := &postgres.StatMetrics{
		Key: pathKey,
	}
	err = s.pg.GetStatRecodByKey(statMetrics)
	if err != nil {
		return Metrics{"0", "0"}, err
	}

	cUA := strconv.Itoa(statMetrics.UniqAgent)
	cReq := strconv.Itoa(statMetrics.Request)

	return Metrics{Request: cReq, UniqAgent: cUA}, nil
}

func generatePathKeyMetrics(urlID string, method string) string {
	key := PathKeyMetrics
	key = strings.Replace(key, "<urlID>", urlID, -1)
	key = strings.Replace(key, "<method>", method, -1)
	return key
}

func generateUAKeyMetrics(urlID string, method string, ua string) string {
	key := UAKeyMetrics
	key = strings.Replace(key, "<urlID>", urlID, -1)
	key = strings.Replace(key, "<method>", method, -1)
	key = strings.Replace(key, "<ua>", ua, -1)
	return key
}
