package logsmonitoring

import "time"

type MonitoringLogs struct {
	ID        int    `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type MonitoringLogsQuery struct {
	UserID    string     `query:"user_id"`
	CreatedAt string     `query:"created_at"`
	DateFrom  *time.Time `query:"date_from"`
	DateTo    *time.Time `query:"date_to"`
	Page      int        `query:"page"`
	PerPage   int        `query:"per_page"`
}

type MonitoringLogsResult struct {
	TotalCount int               `json:"total_count"`
	Logs       []*MonitoringLogs `json:"logs"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
}
