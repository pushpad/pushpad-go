package notification

import "time"

// NotificationAction represents a notification action button.
type NotificationAction struct {
	Title     *string `json:"title,omitempty"`
	TargetURL *string `json:"target_url,omitempty"`
	Icon      *string `json:"icon,omitempty"`
	Action    *string `json:"action,omitempty"`
}

// Notification represents a Pushpad notification.
type Notification struct {
	ID                 int                  `json:"id,omitempty"`
	ProjectID          int                  `json:"project_id,omitempty"`
	Title              string               `json:"title,omitempty"`
	Body               string               `json:"body,omitempty"`
	TargetURL          string               `json:"target_url,omitempty"`
	IconURL            string               `json:"icon_url,omitempty"`
	BadgeURL           string               `json:"badge_url,omitempty"`
	ImageURL           string               `json:"image_url,omitempty"`
	TTL                *int                 `json:"ttl,omitempty"`
	RequireInteraction *bool                `json:"require_interaction,omitempty"`
	Silent             *bool                `json:"silent,omitempty"`
	Urgent             *bool                `json:"urgent,omitempty"`
	CustomData         string               `json:"custom_data,omitempty"`
	Actions            []NotificationAction `json:"actions,omitempty"`
	Starred            *bool                `json:"starred,omitempty"`
	SendAt             *time.Time           `json:"send_at,omitempty"`
	CustomMetrics      []string             `json:"custom_metrics,omitempty"`
	UIDs               []string             `json:"uids"`
	Tags               []string             `json:"tags"`
	CreatedAt          *time.Time           `json:"created_at,omitempty"`
	SuccessfullySent   *int                 `json:"successfully_sent_count,omitempty"`
	OpenedCount        *int                 `json:"opened_count,omitempty"`
	ScheduledCount     *int                 `json:"scheduled_count,omitempty"`
	Scheduled          *bool                `json:"scheduled,omitempty"`
	Cancelled          *bool                `json:"cancelled,omitempty"`
}

// NotificationCreateParams represents a notification create payload.
type NotificationCreateParams struct {
	ProjectID          *int                  `json:"-"`
	Title              *string               `json:"title,omitempty"`
	Body               *string               `json:"body,omitempty"`
	TargetURL          *string               `json:"target_url,omitempty"`
	IconURL            *string               `json:"icon_url,omitempty"`
	BadgeURL           *string               `json:"badge_url,omitempty"`
	ImageURL           *string               `json:"image_url,omitempty"`
	TTL                *int                  `json:"ttl,omitempty"`
	RequireInteraction *bool                 `json:"require_interaction,omitempty"`
	Silent             *bool                 `json:"silent,omitempty"`
	Urgent             *bool                 `json:"urgent,omitempty"`
	CustomData         *string               `json:"custom_data,omitempty"`
	Actions            *[]NotificationAction `json:"actions,omitempty"`
	Starred            *bool                 `json:"starred,omitempty"`
	SendAt             *time.Time            `json:"send_at,omitempty"`
	CustomMetrics      *[]string             `json:"custom_metrics,omitempty"`
	UIDs               *[]string             `json:"uids"`
	Tags               *[]string             `json:"tags"`
}

// NotificationCreateResponse describes the response to creating a notification.
type NotificationCreateResponse struct {
	ID        int        `json:"id"`
	Scheduled *int       `json:"scheduled,omitempty"`
	UIDs      []string   `json:"uids,omitempty"`
	SendAt    *time.Time `json:"send_at,omitempty"`
}

// NotificationListParams controls notification listing.
type NotificationListParams struct {
	ProjectID *int
	Page      *int
}

// NotificationGetParams controls notification fetches.
type NotificationGetParams struct{}

// NotificationCancelParams controls notification cancels.
type NotificationCancelParams struct{}
