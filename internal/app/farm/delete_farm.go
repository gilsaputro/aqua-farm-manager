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

// DeleteFarmRequest is list request parameter for Delete Api
type DeleteFarmRequest struct {
	FarmID   uint   `json:"id"`
	FarmName string `json:"name"`
}

// DeleteFarmResponse is list response parameter for Delete Api
type DeleteFarmResponse struct {
	FarmID   uint   `json:"id"`
	FarmName string `json:"name"`
}

// DeleteFarmHandler is func handler for Delete Farm data
func (h *FarmHandler) DeleteFarmHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[DeleteFarmHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body DeleteFarmRequest
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
	if len(body.FarmName) < 1 && body.FarmID < 1 {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	if len(body.FarmName) > 1 && body.FarmID > 0 {
		code = http.StatusBadRequest
		err = fmt.Errorf("Please Choose to delete by ID or Name")
		return
	}

	errChan := make(chan error, 1)
	var res farm.DeleteDomainResponse
	go func(ctx context.Context) {
		res, err = h.domain.DeleteFarmInfo(farm.DeleteDomainRequest{
			Name: body.FarmName,
			ID:   body.FarmID,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		if err != nil {
			if err == farm.ErrInvalidFarm {
				code = http.StatusNotFound
			} else if err == farm.ErrExistsPonds {
				code = http.StatusConflict
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseDelete(res)
}

func mapResonseDelete(r farm.DeleteDomainResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := DeleteFarmResponse{
		FarmID:   r.ID,
		FarmName: r.Name,
	}
	res.Data = data
	return res
}
