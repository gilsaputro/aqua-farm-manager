package utilhttp

import "net/http"

// WriteResponse is func to generate response for http handler
func WriteResponse(w http.ResponseWriter, data []byte, status int) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return w.Write(data)
}
