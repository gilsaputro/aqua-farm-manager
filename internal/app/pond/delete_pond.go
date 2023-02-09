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

// DeletePondRequest is list request parameter for Delete Api
type DeletePondRequest struct {
	PondID   uint   `json:"id"`
	PondName string `json:"name"`
}

// DeletePondResponse is list response parameter for Delete Api
type DeletePondResponse struct {
	PondID   uint   `json:"id"`
	PondName string `json:"name"`
}

// DeletePondHandler is func handler for Delete Pond data
func (h *PondHandler) DeletePondHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[DeletePondHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body DeletePondRequest
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
	if len(body.PondName) < 1 && body.PondID < 1 {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	if len(body.PondName) > 1 && body.PondID > 0 {
		code = http.StatusBadRequest
		err = fmt.Errorf("Please Choose to delete by ID or Name")
		return
	}

	errChan := make(chan error, 1)
	var res pond.DeleteDomainResponse
	go func(ctx context.Context) {
		res, err = h.domain.DeletePondInfo(pond.DeleteDomainRequest{
			Name: body.PondName,
			ID:   body.PondID,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		if err != nil {
			if err == pond.ErrInvalidPond {
				code = http.StatusOK
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseDelete(res)
}

func mapResonseDelete(r pond.DeleteDomainResponse) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := DeletePondResponse{
		PondID:   r.ID,
		PondName: r.Name,
	}
	res.Data = data
	return res
}
