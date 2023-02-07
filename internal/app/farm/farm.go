package farm

// FarmHandler list dependencies for farm handler
type FarmHandler struct {
	timeoutInSec int
}

// Option set options for http handler config
type Option func(*FarmHandler)

const (
	defaultTimeout = 5
)

// NewFarmHandler is func to create http farm handler
func NewFarmHandler(options ...Option) *FarmHandler {
	handler := &FarmHandler{}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// WithTimeoutOptions is func to set timeout config into handler
func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(fh *FarmHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			fh.timeoutInSec = timeoutinsec
		})
}
