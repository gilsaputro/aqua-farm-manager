package cache

import "strings"

const (
	UniqUAKey    = "count_uniqua"
	RequestedKey = "count_requested"

	TemplateUniqUA          = "p:<path>:<method>:ua:<hash>"
	TemplateTrackingRequest = "p:<path>:<method>"
)

func GetTrackingKey(path, method string) string {
	key := TemplateTrackingRequest
	key = strings.Replace(key, "<path>", path, -1)
	key = strings.Replace(key, "<method>", method, -1)
	return key
}

func GetUniqUAkey(path, method, hash string) string {
	key := TemplateUniqUA
	key = GetTrackingKey(path, method)
	key = strings.Replace(key, "<hash>", hash, -1)
	return key
}
