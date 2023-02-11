package farm

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"aqua-farm-manager/internal/domain/farm"
	utilhttp "aqua-farm-manager/pkg/utilhttp"
)

// UpdateFarmRequest is list request parameter for Update Api
type UpdateFarmRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Owner    string `json:"owner"`
	Area     string `json:"area"`
}

// UpdateFarmResponse is list response parameter for Update Api
type UpdateFarmResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Owner    string `json:"owner"`
	Area     string `json:"area"`
}

// UpdateFarmHandler is func handler for Update Farm data
func (h *FarmHandler) UpdateFarmHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[UpdateFarmHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body UpdateFarmRequest
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
	if len(body.Name) < 1 || (len(body.Location) < 1 && len(body.Area) < 1 && len(body.Owner) < 1) {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	errChan := make(chan error, 1)
	var res farm.UpdateDomainResponse
	go func(ctx context.Context) {
		res, err = h.domain.UpdateFarmInfo(farm.UpdateDomainRequest{
			Name:     body.Name,
			Location: body.Location,
			Owner:    body.Owner,
			Area:     body.Area,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		code = http.StatusGatewayTimeout
		err = fmt.Errorf("Timeout")
		return
	case err = <-errChan:
		if err != nil {
			code = http.StatusInternalServerError
			return
		}
	}

	response = mapResonseUpdate(res)
}

func mapResonseUpdate(r farm.UpdateDomainResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := UpdateFarmResponse{
		ID:       r.ID,
		Name:     r.Name,
		Location: r.Location,
		Owner:    r.Owner,
		Area:     r.Area,
	}
	res.Data = data
	return res
}
