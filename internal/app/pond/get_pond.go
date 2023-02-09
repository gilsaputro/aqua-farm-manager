package pond

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"aqua-farm-manager/internal/domain/pond"
	utilhttp "aqua-farm-manager/pkg/utilhttp"
)

// GetPondRequest is list response parameter for Get Api
type GetPondRequest struct {
	Size   int `json:"size"`
	Cursor int `json:"cursor"`
}

// GetPondResponse is list response parameter for Get Api
type GetPondResponse struct {
	Ponds  []PondInfo `json:"ponds"`
	Cursor *int       `json:"cursor,omitempty"`
}

type PondInfo struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Capacity     float64 `json:"capacity"`
	Depth        float64 `json:"depth"`
	WaterQuality float64 `json:"water_quality"`
	Species      string  `json:"species"`
	FarmID       uint    `json:"farm_id"`
}

// GetPondHandler is func handler for create Pond data
func (h *PondHandler) GetPondHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.timeoutInSec)*time.Second)
	defer cancel()

	var err error
	var response utilhttp.StandardResponse
	var code int = http.StatusOK

	defer func() {
		response.Code = code
		if err == nil {
			response.Message = "success"
		} else {
			response.Message = err.Error()
		}

		data, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			log.Println("[GetPondHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	// checking valid body
	var body GetPondRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Bad Request")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Bad Request")
		return
	}

	if body.Size < 1 || body.Size > 20 {
		body.Size = 20
	}

	if body.Cursor < 1 {
		body.Cursor = 1
	}

	errChan := make(chan error, 1)
	var res []pond.GetPondInfoResponse
	var next int
	go func(ctx context.Context) {
		res, next, err = h.domain.GetAllPond(body.Size, body.Cursor)
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	if len(res) == 0 {
		code = http.StatusNotFound
		err = fmt.Errorf("Data Not Found")
		return
	}

	response = mapResonseGet(res, next)
}

func mapResonseGet(ponds []pond.GetPondInfoResponse, next int) utilhttp.StandardResponse {
	var list []PondInfo

	for _, pond := range ponds {
		info := PondInfo{
			ID:           pond.ID,
			Name:         pond.Name,
			Capacity:     pond.Capacity,
			Depth:        pond.Depth,
			WaterQuality: pond.WaterQuality,
			Species:      pond.Species,
			FarmID:       pond.FarmID,
		}

		list = append(list, info)
	}

	response := GetPondResponse{
		Ponds: list,
	}

	if next > 0 {
		response.Cursor = &next
	}

	return utilhttp.StandardResponse{
		Data: response,
	}
}
