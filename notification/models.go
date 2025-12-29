package notification

import "time"

// NotificationAction represents a notification action button in responses.
type NotificationAction struct {
	Title     string `json:"title"`
	TargetURL string `json:"target_url"`
	Icon      string `json:"icon"`
	Action    string `json:"action"`
}

// Notification represents a Pushpad notification.
type Notification struct {
	ID                 int64                `json:"id"`
	ProjectID          int64                `json:"project_id"`
	Title              string               `json:"title"`
	Body               string               `json:"body"`
	TargetURL          string               `json:"target_url"`
	IconURL            string               `json:"icon_url"`
	BadgeURL           string               `json:"badge_url"`
	ImageURL           string               `json:"image_url"`
	TTL                int64                `json:"ttl"`
	RequireInteraction bool                 `json:"require_interaction"`
	Silent             bool                 `json:"silent"`
	Urgent             bool                 `json:"urgent"`
	CustomData         string               `json:"custom_data"`
	Actions            []NotificationAction `json:"actions"`
	Starred            bool                 `json:"starred"`
	SendAt             time.Time            `json:"send_at"`
	CustomMetrics      []string             `json:"custom_metrics"`
	UIDs               []string             `json:"uids"`
	Tags               []string             `json:"tags"`
	CreatedAt          time.Time            `json:"created_at"`
	SuccessfullySent   int64                `json:"successfully_sent_count"`
	OpenedCount        int64                `json:"opened_count"`
	ScheduledCount     int64                `json:"scheduled_count"`
	Scheduled          bool                 `json:"scheduled"`
	Cancelled          bool                 `json:"cancelled"`
}

// NotificationCreateResponse describes the response to creating a notification.
type NotificationCreateResponse struct {
	ID        int64      `json:"id"`
	Scheduled int64      `json:"scheduled"`
	UIDs      []string   `json:"uids"`
	SendAt    time.Time  `json:"send_at"`
}

// NotificationActionParams represents a notification action button in create payloads.
type NotificationActionParams struct {
	Title     *string `json:"title,omitempty"`
	TargetURL *string `json:"target_url,omitempty"`
	Icon      *string `json:"icon,omitempty"`
	Action    *string `json:"action,omitempty"`
}

// NotificationCreateParams represents a notification create payload.
type NotificationCreateParams struct {
	ProjectID          *int64                `json:"-"`
	Title              *string               `json:"title,omitempty"`
	Body               *string               `json:"body,omitempty"`
	TargetURL          *string               `json:"target_url,omitempty"`
	IconURL            *string               `json:"icon_url,omitempty"`
	BadgeURL           *string               `json:"badge_url,omitempty"`
	ImageURL           *string               `json:"image_url,omitempty"`
	TTL                *int64                `json:"ttl,omitempty"`
	RequireInteraction *bool                 `json:"require_interaction,omitempty"`
	Silent             *bool                 `json:"silent,omitempty"`
	Urgent             *bool                 `json:"urgent,omitempty"`
	CustomData         *string               `json:"custom_data,omitempty"`
	Actions            *[]NotificationActionParams `json:"actions,omitempty"`
	Starred            *bool                 `json:"starred,omitempty"`
	SendAt             *time.Time            `json:"send_at,omitempty"`
	CustomMetrics      *[]string             `json:"custom_metrics,omitempty"`
	UIDs               *[]string             `json:"uids,omitempty"`
	Tags               *[]string             `json:"tags,omitempty"`
}

// NotificationListParams controls notification listing.
type NotificationListParams struct {
	ProjectID *int64
	Page      *int64
}

// NotificationGetParams controls notification fetches.
type NotificationGetParams struct{}

// NotificationCancelParams controls notification cancels.
type NotificationCancelParams struct{}
