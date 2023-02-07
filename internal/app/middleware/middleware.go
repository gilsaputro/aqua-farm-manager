package middleware

import (
	"net/http"
	"sync"

	"aqua-farm-manager/internal/domain/stat"
)

// Middleware struct is list dependecies to run Middleware func
type Middleware struct {
	stat stat.StatDomain
}

// NewMiddleware is func to create Middleware Struct
func NewMiddleware(stat stat.StatDomain) Middleware {
	return Middleware{
		stat: stat,
	}
}

// Middleware is func to validate before execute the handler
func (m *Middleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		ua := r.UserAgent()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			m.stat.IngestStatAPI(path, method, ua)
		}()

		next.ServeHTTP(w, r)
		wg.Wait()
	}
}
