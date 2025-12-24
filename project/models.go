package project

import "time"

// Project represents a Pushpad project.
type Project struct {
	ID                           int64      `json:"id,omitempty"`
	SenderID                     int64      `json:"sender_id,omitempty"`
	Name                         string     `json:"name,omitempty"`
	Website                      string     `json:"website,omitempty"`
	IconURL                      string     `json:"icon_url,omitempty"`
	BadgeURL                     string     `json:"badge_url,omitempty"`
	NotificationsTTL             *int64     `json:"notifications_ttl,omitempty"`
	NotificationsRequireInteract *bool      `json:"notifications_require_interaction,omitempty"`
	NotificationsSilent          *bool      `json:"notifications_silent,omitempty"`
	CreatedAt                    *time.Time `json:"created_at,omitempty"`
}

// ProjectCreateParams is the payload to create a project.
type ProjectCreateParams struct {
	SenderID                     *int64  `json:"sender_id"`
	Name                         *string `json:"name"`
	Website                      *string `json:"website"`
	IconURL                      *string `json:"icon_url,omitempty"`
	BadgeURL                     *string `json:"badge_url,omitempty"`
	NotificationsTTL             *int64  `json:"notifications_ttl,omitempty"`
	NotificationsRequireInteract *bool   `json:"notifications_require_interaction,omitempty"`
	NotificationsSilent          *bool   `json:"notifications_silent,omitempty"`
}

// ProjectUpdateParams is the payload to update a project.
type ProjectUpdateParams struct {
	Name                         *string `json:"name,omitempty"`
	Website                      *string `json:"website,omitempty"`
	IconURL                      *string `json:"icon_url,omitempty"`
	BadgeURL                     *string `json:"badge_url,omitempty"`
	NotificationsTTL             *int64  `json:"notifications_ttl,omitempty"`
	NotificationsRequireInteract *bool   `json:"notifications_require_interaction,omitempty"`
	NotificationsSilent          *bool   `json:"notifications_silent,omitempty"`
}

// ProjectListParams controls project listing.
type ProjectListParams struct{}

// ProjectGetParams controls project fetches.
type ProjectGetParams struct{}

// ProjectDeleteParams controls project deletes.
type ProjectDeleteParams struct{}
