package logsmonitoring

import "context"

type Service interface {
	MonotoringLogs(ctx context.Context, query *MonitoringLogsQuery) (*MonitoringLogsResult, error)
}
