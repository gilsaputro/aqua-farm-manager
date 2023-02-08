package pond

// PondHandler list dependencies for pond handler
type PondHandler struct {
	timeoutInSec int
}

// Option set options for http handler config
type Option func(*PondHandler)

const (
	defaultTimeout = 5
)

// NewPondHandler is func to create http Pond handler
func NewPondHandler(options ...Option) *PondHandler {
	handler := &PondHandler{}

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
