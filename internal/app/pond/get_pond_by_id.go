package pond

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"aqua-farm-manager/internal/domain/pond"
	utilhttp "aqua-farm-manager/pkg/utilhttp"

	"github.com/gorilla/mux"
)

// GetByIDPondResponse is list response parameter for GetByID Api
type GetByIDPondResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Capacity     float64   `json:"capacity"`
	Depth        float64   `json:"depth"`
	WaterQuality float64   `json:"water_quality"`
	Species      string    `json:"species"`
	FarmInfo     *FarmInfo `json:"farm,omitempty"`
}

// FarmInfo is list parameter for farm info
type FarmInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Owner    string `json:"owner"`
	Area     string `json:"area"`
}

// GetByIDPondHandler is func handler for create Pond data
func (h *PondHandler) GetByIDPondHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[GetByIDPondHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	// checking valid body
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	errChan := make(chan error, 1)
	var res pond.GetPondInfoResponse
	go func(ctx context.Context) {
		res, err = h.domain.GetPondInfoByID(uint(id))
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		code = http.StatusGatewayTimeout
		err = fmt.Errorf("Timeout")
		return
	case err = <-errChan:
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				code = http.StatusNotFound
				err = fmt.Errorf("Data Not Found")
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseGetByID(res)
}

func mapResonseGetByID(r pond.GetPondInfoResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse

	data := GetByIDPondResponse{
		ID:           r.ID,
		Name:         r.Name,
		Capacity:     r.Capacity,
		Depth:        r.Depth,
		WaterQuality: r.WaterQuality,
		Species:      r.Species,
	}
	if r.FarmInfo.ID != 0 {
		data.FarmInfo = &FarmInfo{
			ID:       r.FarmInfo.ID,
			Name:     r.FarmInfo.Name,
			Location: r.FarmInfo.Location,
			Owner:    r.FarmInfo.Owner,
			Area:     r.FarmInfo.Area,
		}
	}

	res.Data = data
	return res
}
