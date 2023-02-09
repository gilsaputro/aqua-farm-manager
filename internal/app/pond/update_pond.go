package pond

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"aqua-farm-manager/internal/domain/farm"
	"aqua-farm-manager/internal/domain/pond"
	utilhttp "aqua-farm-manager/pkg/utilhttp"
)

// UpdatePondRequest is list request parameter for Update Api
type UpdatePondRequest struct {
	Name         string  `json:"name"`
	Capacity     float64 `json:"capacity"`
	Depth        float64 `json:"depth"`
	WaterQuality float64 `json:"water_quality"`
	Species      string  `json:"species"`
	FarmID       uint    `json:"farm_id"`
}

// UpdatePondResponse is list response parameter for Update Api
type UpdatePondResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Capacity     float64 `json:"capacity"`
	Depth        float64 `json:"depth"`
	WaterQuality float64 `json:"water_quality"`
	Species      string  `json:"species"`
	FarmID       uint    `json:"farm_id"`
}

// UpdatePondHandler is func handler for Update Pond data
func (h *PondHandler) UpdatePondHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[UpdatePondHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body UpdatePondRequest
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

	// checking valid body
	if len(body.Name) < 1 || (len(body.Species) < 1 && body.Capacity < 1 && body.Depth < 1 && body.WaterQuality < 1 && body.FarmID < 1) {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	errChan := make(chan error, 1)
	var res pond.UpdateDomainResponse
	go func(ctx context.Context) {
		res, err = h.domain.UpdatePondInfo(pond.UpdateDomainRequest{
			Name:         body.Name,
			Capacity:     body.Capacity,
			Depth:        body.Capacity,
			WaterQuality: body.WaterQuality,
			Species:      body.Species,
			FarmID:       body.FarmID,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		if err != nil {
			if err == farm.ErrDuplicateFarm {
				code = http.StatusBadRequest
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseUpdate(res)
}

func mapResonseUpdate(r pond.UpdateDomainResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := UpdatePondResponse{
		ID:           r.ID,
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
		FarmID:       r.FarmID,
	}
	res.Data = data
	return res
}
