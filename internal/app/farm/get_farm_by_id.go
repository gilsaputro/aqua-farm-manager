package farm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"aqua-farm-manager/internal/domain/farm"
	utilhttp "aqua-farm-manager/pkg/utilhttp"

	"github.com/gorilla/mux"
)

// FFarmResponse is list response parameter for GetByID Api
type GetByIDFarmResponse struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Location string      `json:"location"`
	Owner    string      `json:"owner"`
	Area     string      `json:"area"`
	PondInfo *[]PondInfo `json:"pond_info,omitempty"`
}

type PondInfo struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Capacity     float64 `json:"capacity"`
	Depth        float64 `json:"depth"`
	WaterQuality float64 `json:"water_quality"`
	Species      string  `json:"species"`
	Status       int     `json:"status"`
}

// GetByIDFarmHandler is func handler for create Farm data
func (h *FarmHandler) GetByIDFarmHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[GetByIDFarmHandler]-Error Marshal Response :", err)
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
	var res farm.GetFarmInfoResponse
	go func(ctx context.Context) {
		res, err = h.domain.GetFarmInfoByID(uint(id))
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

func mapResonseGetByID(r farm.GetFarmInfoResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse

	var list []PondInfo
	for _, pond := range r.PondInfos {
		list = append(list, PondInfo{
			ID:           pond.ID,
			Name:         pond.Name,
			Capacity:     pond.Capacity,
			Depth:        pond.Depth,
			WaterQuality: pond.WaterQuality,
			Species:      pond.Species,
		})
	}

	data := GetByIDFarmResponse{
		ID:       r.ID,
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
	}

	if len(list) > 0 {
		data.PondInfo = &list
	}

	res.Data = data
	return res
}
