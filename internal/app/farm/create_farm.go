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

// CreateFarmRequest is list request parameter for Create Api
type CreateFarmRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Owner    string `json:"owner"`
	Area     string `json:"area"`
}

// CreateFarmResponse is list response parameter for Create Api
type CreateFarmResponse struct {
	ID uint `json:"id"`
}

// CreateFarmHandler is func handler for create Farm data
func (h *FarmHandler) CreateFarmHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[CreateFarmHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body CreateFarmRequest
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
	var res farm.CreateDomainResponse
	go func(ctx context.Context) {
		res, err = h.domain.CreateFarmInfo(farm.CreateDomainRequest{
			Name:     body.Name,
			Location: body.Location,
			Owner:    body.Owner,
			Area:     body.Area,
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

	response = mapResonseCreate(res)
}

func mapResonseCreate(r farm.CreateDomainResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := CreateFarmResponse{
		ID: r.ID,
	}
	res.Data = data
	return res
}
