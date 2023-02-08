package stat

// IngestMetricsRequest list is request  for IngestMetrics
type IngestMetricsRequest struct {
	UrlID     string
	Method    string
	UA        string
	IsSuccess bool
}

// GetMetricsRequest list is request  for GetMetrics
type GetMetricsRequest struct {
	UrlID  string
	Method string
}

// BackupMetricsRequest list is request  for BackupMetrics
type BackupMetricsRequest struct {
	UrlID   string
	Method  string
	Metrics MetricsRequest
}

// MigrateMetricsRequest list is request  for MigrateMetrics
type MigrateMetricsRequest struct {
	UrlID   string
	Method  string
	Metrics MetricsInfo
}

// GetStatDataRequest list is request  for GetStatData
type GetStatDataRequest struct {
	UrlID  string
	Method string
}

type MetricsRequest struct {
	NumRequest   int
	NumUniqAgent int
	NumSuccess   int
	NumError     int
}

type MetricsInfo struct {
	NumRequest   string
	NumUniqAgent string
	NumSuccess   string
	NumError     string
}
