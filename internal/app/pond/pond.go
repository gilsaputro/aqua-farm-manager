package pond

import "aqua-farm-manager/internal/domain/pond"

// PondHandler list dependencies for pond handler
type PondHandler struct {
	timeoutInSec int
	domain       pond.PondDomain
}

// Option set options for http handler config
type Option func(*PondHandler)

const (
	defaultTimeout = 5
)

// NewPondHandler is func to create http Pond handler
func NewPondHandler(domain pond.PondDomain, options ...Option) *PondHandler {
	handler := &PondHandler{
		domain:       domain,
		timeoutInSec: defaultTimeout,
	}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// WithTimeoutOptions is func to set timeout config into handler
func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(fh *PondHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			fh.timeoutInSec = timeoutinsec
		})
}
