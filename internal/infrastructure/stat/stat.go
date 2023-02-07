package stat

import (
	"strings"
	"sync"

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
	GetMetrics(urlID, method string) (string, string, error)
}

// Stat is list dependencies stat store
type Stat struct {
	redis redis.RedisMethod
}

// NewStatStore is func to generate StatStore interface
func NewStatStore(redis redis.RedisMethod) StatStore {
	return &Stat{
		redis: redis,
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
func (s *Stat) GetMetrics(urlID, method string) (string, string, error) {
	var err error
	pathKey := generatePathKeyMetrics(urlID, method)

	metrics, err := s.redis.HGETALL(pathKey)
	if err != nil {
		return "", "", err
	}

	cUA := metrics[CountUA]
	cReq := metrics[CountRequested]

	return cUA, cReq, nil
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
