package stat

import (
	"aqua-farm-manager/internal/domain/stat"
	"aqua-farm-manager/pkg/utilhttp"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StatHandler struct is list dependecies to run Stat Handler
type StatHandler struct {
	stat         stat.StatDomain
	timeoutInSec int
}

// NewMiddleware is func to create StatHandler Struct
func NewStatHandler(stat stat.StatDomain, options ...Option) *StatHandler {
	handler := &StatHandler{
		stat:         stat,
		timeoutInSec: defaultTimeout,
	}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// Option set options for http handler config
type Option func(*StatHandler)

const (
	defaultTimeout = 5
)

// WithTimeoutOptions is func to set timeout config into handler
func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(h *StatHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}

// GetStatHandler is func handler for Get Stat API
func (h *StatHandler) GetStatHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.timeoutInSec)*time.Second)
	defer cancel()

	var err error
	var response Response
	var code int = 200

	defer func() {
		response.Code = code
		if err == nil {
			response.Message = "success"
		} else {
			response.Message = err.Error()
		}

		data, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			log.Println("[GetStatHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	errChan := make(chan error, 1)
	var metrics map[string]stat.StatMetrics
	go func(ctx context.Context) {
		metrics = h.stat.GenerateStatAPI()
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	response = mapResponse(metrics)
}

func mapResponse(metrics map[string]stat.StatMetrics) Response {
	var res Response
	var mapMetrics = make(map[string]Metrics, len(metrics))
	for key, value := range metrics {
		mapMetrics[key] = Metrics{
			Count:           value.Requested,
			UniqueUserAgent: value.UniqAgent,
		}
	}
	res.Data = mapMetrics
	return res
}

// InitMigrate is func to init scheduler to backup data stat from redis to postgres
func (h *StatHandler) InitMigrate(tickInMinute int) {
	go func() {
		tick := time.Tick(time.Duration(tickInMinute) * time.Minute)
		for range tick {
			log.Println("Running Backup Stat")
			h.stat.BackUpStat()
		}
	}()
}
