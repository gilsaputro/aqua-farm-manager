package stat

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"aqua-farm-manager/internal/model"
	"aqua-farm-manager/pkg/postgres"
	"aqua-farm-manager/pkg/redis"

	"github.com/jinzhu/gorm"
)

var (
	PathKeyMetrics = "P:<urlID>:<method>"
	UAKeyMetrics   = "P:<urlID>:<method>:<ua>"

	CountUA        = "Count_UA"
	CountRequested = "Count_Req"
	CountSuccess   = "Count_Success"
	CountError     = "Count_Error"
)

// StatStore is set of methods for interacting with a metric storage system
type StatStore interface {
	IngestMetrics(IngestMetricsRequest) error
	GetMetrics(GetMetricsRequest) (MetricsInfo, error)
	BackupMetrics(BackupMetricsRequest) error
	MigrateMetrics(MigrateMetricsRequest) error
	GetStatData(GetStatDataRequest) (MetricsInfo, error)
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
func (s *Stat) IngestMetrics(r IngestMetricsRequest) error {
	var err error
	pathKey := generatePathKeyMetrics(r.UrlID, r.Method)
	uakey := generateUAKeyMetrics(r.UrlID, r.Method, r.UA)

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		if r.IsSuccess {
			s.redis.HINCRBY(pathKey, CountSuccess)
		} else {
			s.redis.HINCRBY(pathKey, CountError)
		}
	}()

	wg.Wait()
	return err
}

// GetMetrics is func to get api metrics from redis
func (s *Stat) GetMetrics(r GetMetricsRequest) (MetricsInfo, error) {
	var err error
	pathKey := generatePathKeyMetrics(r.UrlID, r.Method)

	metrics, err := s.redis.HGETALL(pathKey)
	if err != nil {
		return MetricsInfo{"0", "0", "0", "0"}, err
	}

	numUA := metrics[CountUA]
	numReq := metrics[CountRequested]
	numSuc := metrics[CountSuccess]
	numErr := metrics[CountError]

	if len(numUA) == 0 {
		numUA = "0"
	}

	if len(numReq) == 0 {
		numReq = "0"
	}

	if len(numErr) == 0 {
		numErr = "0"
	}

	if len(numSuc) == 0 {
		numSuc = "0"
	}

	return MetricsInfo{
		NumRequest:   numReq,
		NumUniqAgent: numUA,
		NumSuccess:   numSuc,
		NumError:     numErr,
	}, nil
}

// BackupMetrics is func to backup metrics from redis to postgres
func (s *Stat) BackupMetrics(r BackupMetricsRequest) error {
	var err error
	pathKey := generatePathKeyMetrics(r.UrlID, r.Method)
	stat := postgres.StatMetrics{
		Key:        pathKey,
		Request:    r.Metrics.NumRequest,
		UniqAgent:  r.Metrics.NumUniqAgent,
		NumSuccess: r.Metrics.NumSuccess,
		NumError:   r.Metrics.NumError,
	}

	db := s.pg.GetDB()
	if db == nil {
		return errors.New("Database Client is not init")
	}

	err = updateStat(db, &stat)
	if err != nil {
		return err
	}
	stat.Status = model.Active.Value()
	err = insert(db, &stat)

	return err
}

// MigrateMetrics is func to migrate metrics from postgres to redis
func (s *Stat) MigrateMetrics(r MigrateMetricsRequest) error {
	var errUA, errReq, errSuc, errEr error
	key := generatePathKeyMetrics(r.UrlID, r.Method)
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		errUA = s.redis.HSET(key, CountUA, r.Metrics.NumUniqAgent)
		if errUA != nil {
			log.Println("MigrateMetrics-Error Ingest Uniq Agent :", errUA)
		}
	}()

	go func() {
		defer wg.Done()
		errReq = s.redis.HSET(key, CountRequested, r.Metrics.NumRequest)
		if errReq != nil {
			log.Println("MigrateMetrics-Error Ingest Requested :", errReq)
		}
	}()

	go func() {
		defer wg.Done()
		errEr = s.redis.HSET(key, CountError, r.Metrics.NumError)
		if errEr != nil {
			log.Println("MigrateMetrics-Error Ingest Requested :", errReq)
		}
	}()

	go func() {
		defer wg.Done()
		errSuc = s.redis.HSET(key, CountSuccess, r.Metrics.NumSuccess)
		if errSuc != nil {
			log.Println("MigrateMetrics-Error Ingest Requested :", errReq)
		}
	}()

	wg.Wait()

	if errUA != nil || errReq != nil || errEr != nil || errSuc != nil {
		return fmt.Errorf("got error while migrate")
	}

	return nil
}

// GetStatData is func to metrics from postgres
func (s *Stat) GetStatData(r GetStatDataRequest) (MetricsInfo, error) {
	var err error
	pathKey := generatePathKeyMetrics(r.UrlID, r.Method)

	statMetrics := &postgres.StatMetrics{
		Key: pathKey,
	}

	db := s.pg.GetDB()
	if db == nil {
		return MetricsInfo{"0", "0", "0", "0"}, errors.New("Database Client is not init")
	}
	err = getStatRecodByKey(db, statMetrics)
	if err != nil {
		return MetricsInfo{"0", "0", "0", "0"}, err
	}

	cUA := strconv.Itoa(statMetrics.UniqAgent)
	cReq := strconv.Itoa(statMetrics.Request)
	cSuccess := strconv.Itoa(statMetrics.NumSuccess)
	cError := strconv.Itoa(statMetrics.NumError)

	return MetricsInfo{
		NumRequest:   cReq,
		NumUniqAgent: cUA,
		NumSuccess:   cSuccess,
		NumError:     cError,
	}, nil
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

// checkStatExists is func to check if the data is stat exist
func checkStatExists(db *gorm.DB, stat postgres.StatMetrics) bool {
	var count = int64(0)
	db.Model(stat).Where("key = ?", stat.Key).Count(&count).Limit(1)
	return count > 0
}

// insert is func to insert data into table
func insert(db *gorm.DB, stat *postgres.StatMetrics) error {
	return db.Create(stat).Error
}

// getStatRecodByKey is func to get data into stat table using key
func getStatRecodByKey(db *gorm.DB, stat *postgres.StatMetrics) error {
	return db.Where("key = ? and status = ?", stat.Key, model.Active.Value()).First(stat).Error
}

// updateStat is func to update data into stat table
func updateStat(db *gorm.DB, stat *postgres.StatMetrics) error {
	return db.Model(stat).Where("key = ? and status = ?", stat.Key, model.Active.Value()).Update(postgres.StatMetrics{Status: model.Inactive.Value()}).Error
}
