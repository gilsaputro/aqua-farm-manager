package farm

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"aqua-farm-manager/internal/domain/farm"
	utilhttp "aqua-farm-manager/pkg/utilhttp"
)

// GetFarmRequest is list response parameter for Get Api
type GetFarmRequest struct {
	Size   int `json:"size"`
	Cursor int `json:"cursor"`
}

// GetFarmResponse is list response parameter for Get Api
type GetFarmResponse struct {
	Farms  []FarmInfo `json:"farms"`
	Cursor *int       `json:"cursor,omitempty"`
}

type FarmInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Owner    string `json:"owner"`
	Area     string `json:"area"`
	PondIDs  []uint `json:"list_pondID"`
}

// GetFarmHandler is func handler for create Farm data
func (h *FarmHandler) GetFarmHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[GetFarmHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	// checking valid body
	var body GetFarmRequest
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
	var res []farm.GetFarmInfoResponse
	var next int
	go func(ctx context.Context) {
		res, next, err = h.domain.GetFarm(body.Size, body.Cursor)
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

	if len(res) == 0 {
		code = http.StatusNotFound
		err = fmt.Errorf("Data Not Found")
		return
	}

	response = mapResonseGet(res, next)
}

func mapResonseGet(farms []farm.GetFarmInfoResponse, next int) utilhttp.StandardResponse {
	var list []FarmInfo

	for _, farm := range farms {
		info := FarmInfo{
			ID:       farm.ID,
			Name:     farm.Name,
			Location: farm.Location,
			Owner:    farm.Owner,
			Area:     farm.Area,
			PondIDs:  farm.PondIDs,
		}

		list = append(list, info)
	}

	response := GetFarmResponse{
		Farms: list,
	}

	if next > 0 {
		response.Cursor = &next
	}

	return utilhttp.StandardResponse{
		Data: response,
	}
}
