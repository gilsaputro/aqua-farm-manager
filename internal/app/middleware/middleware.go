package middleware

import (
	"aqua-farm-manager/internal/domain/cache"
	"aqua-farm-manager/internal/infrastructure/redis"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/crypto/sha3"
)

// Middleware struct is list dependecies to run Middleware func
type Middleware struct {
	redis redis.RedisMethod
}

// NewMiddleware is func to create Middleware Struct
func NewMiddleware(redis redis.RedisMethod) Middleware {
	return Middleware{
		redis: redis,
	}
}

// Middleware is func to validate before execute the handler
func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		// generate uniq user agent key
		ua := r.UserAgent()
		hash := fmt.Sprintf("%x", sha3.Sum256([]byte(ua)))
		uakey := cache.GetUniqUAkey(path, method, hash)

		// set key to redis, if the key is exists it's not uniq
		isNew, err := m.redis.SETNX(uakey)
		if err != nil {
			fmt.Println(err)
		}

		//generate tracking key
		uarequested := cache.GetTrackingKey(path, method)
		wg := sync.WaitGroup{}
		// incr count uniq ua if is new
		if isNew {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.redis.HINCRBY(uarequested, cache.UniqUAKey)
			}()
		}

		// incr count requested
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.redis.HINCRBY(uarequested, cache.RequestedKey)
		}()

		next.ServeHTTP(w, r)
		wg.Wait()
	})
}
