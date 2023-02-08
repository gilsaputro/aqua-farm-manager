package middleware

import (
	"fmt"
	"net/http"

	"aqua-farm-manager/internal/app/trackingevent"
	"aqua-farm-manager/pkg/nsq"
)

// Middleware struct is list dependecies to run Middleware func
type Middleware struct {
	nsq   nsq.NsqMethod
	topic string
}

// NewMiddleware is func to create Middleware Struct
func NewMiddleware(topic string, nsq nsq.NsqMethod) Middleware {
	return Middleware{
		nsq:   nsq,
		topic: topic,
	}
}

// statusResponseWriter is a custom ResponseWriter type that wraps an existing http.ResponseWriter
// and adds the ability to track the HTTP status code
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader is a custom implementation of the http.ResponseWriter's WriteHeader method
// that tracks the HTTP status code by storing it in the statusCode field.
func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Middleware is func to validate before execute the handler
func (m *Middleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		ua := r.UserAgent()
		sw := &statusResponseWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r)
		go func() {
			m.publishToTrackingEvent(path, method, ua, sw.statusCode)
		}()
	}
}

func (m *Middleware) publishToTrackingEvent(path, method, ua string, code int) {
	msg := trackingevent.TrackingEventMessage{
		Path:   path,
		Code:   code,
		Method: method,
		UA:     ua,
	}

	err := m.nsq.Publish(m.topic, msg)
	if err != nil {
		fmt.Println("Middleware-Got Error while Publish :", err)
	}
}
