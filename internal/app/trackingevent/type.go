package trackingevent

//TrackingEventMessage represents data object of nsq message for aqua_farm_tracking_event
type TrackingEventMessage struct {
	Path   string `json:"path"`
	Code   int    `json:"code"`
	Method string `json:"method"`
	UA     string `json:"ua"`
}
