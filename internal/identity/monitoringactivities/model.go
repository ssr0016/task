package monitoringactivities

import "errors"

type ActivityLog struct {
	ID        int    `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"user_id"`
	Activity  string `db:"activity" json:"activity"`
	Action    string `db:"action" json:"action"`
	Resource  string `db:"resource" json:"resource"`
	Details   string `db:"details" json:"details"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type CreateActivityLogCommand struct {
	UserID    string `json:"user_id"` // Change UserID to string
	Activity  string `json:"activity"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Details   string `json:"details"`
	CreatedAt string `json:"created_at"`
}

type SearchLogActivityQuery struct {
	UserID    string `query:"user_id"`
	Activity  string `query:"activity"`
	Action    string `query:"action"`
	Resource  string `query:"resource"`
	Details   string `query:"details"`
	CreatedAt string `query:"created_at"`
	Page      int    `query:"page"`
	PerPage   int    `query:"per_page"`
}

type SearchLogActivityResult struct {
	TotalCount int            `json:"total_count"`
	Activities []*ActivityLog `json:"activities"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}

func (cmd *CreateActivityLogCommand) Validate() error {
	if cmd.UserID == "" || cmd.Activity == "" || cmd.Action == "" || cmd.Resource == "" {
		return errors.New("user_id, activity, action, and resource fields are required")
	}

	return nil
}
