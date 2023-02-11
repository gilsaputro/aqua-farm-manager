package stat

import (
	"aqua-farm-manager/internal/app"
	"aqua-farm-manager/internal/infrastructure/stat"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"golang.org/x/crypto/sha3"
)

// StatDomain is list method for stat domain
type StatDomain interface {
	GenerateStatAPI() map[string]StatMetrics
	IngestStatAPI(IngestStatRequest)
	BackUpStat()
	MigrateStat()
}

type IngestStatRequest struct {
	Path   string
	Method string
	Ua     string
	Code   int
}

// Stat is list dependencies stat domain
type Stat struct {
	store stat.StatStore
}

// StatMetrics denotes list stat value of api
type StatMetrics struct {
	NumUniqAgent int
	NumRequested int
	NumSuccess   int
	NumError     int
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
			metric, err := s.store.GetMetrics(stat.GetMetricsRequest{
				UrlID:  url,
				Method: method,
			})
			if err != nil {
				// get data using database
				metric, err = s.store.GetStatData(stat.GetStatDataRequest{
					UrlID:  url,
					Method: method,
				})
				if err != nil {
					continue
				}
			}

			count_req, _ := strconv.Atoi(metric.NumRequest)
			count_ua, _ := strconv.Atoi(metric.NumUniqAgent)
			count_suc, _ := strconv.Atoi(metric.NumSuccess)
			count_err, _ := strconv.Atoi(metric.NumError)
			if count_req != 0 || count_ua != 0 {
				key := method + " " + id.String()
				metrics[key] = StatMetrics{
					NumUniqAgent: count_ua,
					NumRequested: count_req,
					NumSuccess:   count_suc,
					NumError:     count_err,
				}
			}
		}
	}
	return metrics
}

// IngestStatAPI is func to ingest stat metrics based on path and method
func (s *Stat) IngestStatAPI(r IngestStatRequest) {
	urlID := app.UrlIDValue[r.Path]

	if urlID.Int() != 0 {
		hash := fmt.Sprintf("%x", sha3.Sum256([]byte(r.Ua)))
		url := strconv.Itoa(urlID.Int())
		err := s.store.IngestMetrics(
			stat.IngestMetricsRequest{
				UrlID:     url,
				Method:    r.Method,
				UA:        hash,
				IsSuccess: r.Code == http.StatusOK,
			},
		)
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
			metric, err := s.store.GetMetrics(stat.GetMetricsRequest{UrlID: url, Method: method})
			if err != nil {
				continue
			}

			count_req, _ := strconv.Atoi(metric.NumRequest)
			count_ua, _ := strconv.Atoi(metric.NumUniqAgent)
			count_suc, _ := strconv.Atoi(metric.NumSuccess)
			count_err, _ := strconv.Atoi(metric.NumError)

			err = s.store.BackupMetrics(stat.BackupMetricsRequest{
				UrlID:  url,
				Method: method,
				Metrics: stat.MetricsRequest{
					NumRequest:   count_req,
					NumUniqAgent: count_ua,
					NumSuccess:   count_suc,
					NumError:     count_err,
				},
			})
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
				metric, err := s.store.GetStatData(stat.GetStatDataRequest{
					UrlID:  url,
					Method: method,
				})
				if err != nil {
					return
				}

				err = s.store.MigrateMetrics(
					stat.MigrateMetricsRequest{
						UrlID:   url,
						Method:  method,
						Metrics: metric,
					},
				)
				if err != nil {
					fmt.Println("[MigrateStat]-Got Error:", err)
					return
				}
			}(url, method)
		}
	}
	wg.Wait()
}
