package stat

import (
	"aqua-farm-manager/internal/app"
	"aqua-farm-manager/internal/infrastructure/stat"
	"fmt"
	"strconv"
	"sync"

	"golang.org/x/crypto/sha3"
)

// StatDomain is list method for stat domain
type StatDomain interface {
	GenerateStatAPI() map[string]StatMetrics
	IngestStatAPI(path, method, ua string)
	BackUpStat()
	MigrateStat()
}

// Stat is list dependencies stat domain
type Stat struct {
	store stat.StatStore
}

// StatMetrics denotes list stat value of api
type StatMetrics struct {
	UniqAgent int
	Requested int
}

// NewStatDomain is func to generat StatDomain interface
func NewStatDomain(store stat.StatStore) StatDomain {
	return &Stat{
		store: store,
	}
}

// GenerateStatAPI is func to generate stat info for all api
func (s *Stat) GenerateStatAPI() map[string]StatMetrics {
	var metrics = make(map[string]StatMetrics, app.Limit-1)
	for id := app.UrlID(1); id < app.Limit; id++ {
		url := strconv.Itoa(id.Int())
		listmethod := app.UrlIDMethod[id]
		for _, method := range listmethod {
			metric, err := s.store.GetMetrics(url, method)
			if err != nil {
				continue
			}

			count_req, _ := strconv.Atoi(metric.Request)
			count_ua, _ := strconv.Atoi(metric.UniqAgent)
			if count_req != 0 || count_ua != 0 {
				key := method + " " + id.String()
				metrics[key] = StatMetrics{
					UniqAgent: count_ua,
					Requested: count_req,
				}
			}
		}
	}
	return metrics
}

// IngestStatAPI is func to ingest stat metrics based on path and method
func (s *Stat) IngestStatAPI(path, method, ua string) {
	urlID := app.UrlIDValue[path]

	if urlID.Int() != 0 {
		hash := fmt.Sprintf("%x", sha3.Sum256([]byte(ua)))
		url := strconv.Itoa(urlID.Int())
		err := s.store.IngestMetrics(url, method, hash)
		if err != nil {
			fmt.Println("[IngestStatAPI]-Got Error:", err)
		}
	}
}

// BackUpStat is func to backup data from redis to postgres
func (s *Stat) BackUpStat() {
	for id := app.UrlID(1); id < app.Limit; id++ {
		url := strconv.Itoa(id.Int())
		listmethod := app.UrlIDMethod[id]
		for _, method := range listmethod {
			metric, err := s.store.GetMetrics(url, method)
			if err != nil {
				continue
			}

			count_req, _ := strconv.Atoi(metric.Request)
			count_ua, _ := strconv.Atoi(metric.UniqAgent)

			err = s.store.BackupMetrics(url, method, count_req, count_ua)
			if err != nil {
				fmt.Println("[BackUpStat]-Got Error:", err)
				continue
			}
		}
	}
}

// MigrateStat is func to migrate data from postgres to redis
func (s *Stat) MigrateStat() {
	var wg sync.WaitGroup
	for id := app.UrlID(1); id < app.Limit; id++ {
		url := strconv.Itoa(id.Int())
		listmethod := app.UrlIDMethod[id]
		for _, method := range listmethod {
			wg.Add(1)
			go func(url, method string) {
				defer wg.Done()
				metric, err := s.store.GetStatData(url, method)
				if err != nil {
					return
				}

				err = s.store.MigrateMetrics(url, method, metric.Request, metric.UniqAgent)
				if err != nil {
					fmt.Println("[MigrateStat]-Got Error:", err)
					return
				}
			}(url, method)
		}
	}
	wg.Wait()
}
