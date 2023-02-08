package pond

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"aqua-farm-manager/internal/domain/pond"
	utilhttp "aqua-farm-manager/pkg/utilhttp"
)

// CreatePondRequest is list request parameter for Create Api
type CreatePondRequest struct {
	Name         string  `json:"name"`
	Capacity     float64 `json:"capacity"`
	Depth        float64 `json:"depth"`
	WaterQuality float64 `json:"water_quality"`
	Species      string  `json:"species"`
	FarmID       uint    `json:"farm_id"`
}

// CreatePondResponse is list response parameter for Create Api
type CreatePondResponse struct {
	PondID uint `json:"pond_id"`
}

// CreatePondHandler is func handler for create pond data
func (h *PondHandler) CreatePondHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[CreatePondHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body CreatePondRequest
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
	if len(body.Name) < 1 || body.FarmID < 1 || (len(body.Species) < 1 && body.Capacity < 1 && body.Depth < 1 && body.WaterQuality < 1) {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	errChan := make(chan error, 1)
	var res pond.PondResponse
	go func(ctx context.Context) {
		res, err = h.domain.CreatePondInfo(pond.PondDomainRequest{
			Name:         body.Name,
			Capacity:     body.Capacity,
			Depth:        body.Depth,
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
			if err == pond.ErrDuplicatePond || err == pond.ErrInvalidFarm {
				code = http.StatusBadRequest
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonse(res)
}

func mapResonse(r pond.PondResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := CreatePondResponse{
		PondID: r.PondID,
	}
	res.Data = data
	return res
}
