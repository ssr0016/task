package monitoringactivities

import (
	"context"
)

type Service interface {
	LogActivity(ctx context.Context, cmd *CreateActivityLogCommand) error
	SearchLogActivities(ctx context.Context, query *SearchLogActivityQuery) (*SearchLogActivityResult, error)
}
