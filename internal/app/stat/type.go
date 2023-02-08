package stat

type Response struct {
	Data    *map[string]Metrics `json:"data,omitempty"`
	Code    int                 `json:"code"`
	Message string              `json:"message"`
}

type Metrics struct {
	Count           int `json:"count"`
	UniqueUserAgent int `json:"unique_user_agent"`
	NumSuccess      int `json:"num_success"`
	NumError        int `json:"num_error"`
}
