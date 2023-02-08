package pond

import (
	"context"
	"fmt"
	"net/http"
	"time"

	utilhttp "aqua-farm-manager/pkg/utilhttp"
)

// GetFarmHandler is func handler for get pond data
func (h *PondHandler) GetFarmHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.timeoutInSec)*time.Second)
	defer cancel()

	var err error
	var response []byte
	var code int = 200

	defer func() {
		utilhttp.WriteResponse(w, response, code)
	}()

	errChan := make(chan error, 1)

	go func(ctx context.Context) {
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	response = []byte(`{"hello":"GET"}`)
}
