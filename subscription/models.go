package subscription

import "time"

// Subscription represents a Pushpad subscription.
type Subscription struct {
	ID          int64      `json:"id,omitempty"`
	ProjectID   int64      `json:"project_id,omitempty"`
	Endpoint    string     `json:"endpoint,omitempty"`
	P256DH      string     `json:"p256dh,omitempty"`
	Auth        string     `json:"auth,omitempty"`
	UID         string     `json:"uid,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	LastClickAt time.Time  `json:"last_click_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
}

// SubscriptionCreateParams is the payload to create a subscription.
type SubscriptionCreateParams struct {
	ProjectID *int64    `json:"-"`
	Endpoint  *string   `json:"endpoint"`
	P256DH    *string   `json:"p256dh,omitempty"`
	Auth      *string   `json:"auth,omitempty"`
	UID       *string   `json:"uid,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
}

// SubscriptionUpdateParams is the payload to update a subscription.
type SubscriptionUpdateParams struct {
	ProjectID *int64    `json:"-"`
	UID       *string   `json:"uid,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
}

// SubscriptionListParams controls subscription listing.
type SubscriptionListParams struct {
	ProjectID *int64
	Page      *int64
	PerPage   *int64
	UIDs      *[]string
	Tags      *[]string
}

// SubscriptionGetParams controls subscription fetches.
type SubscriptionGetParams struct {
	ProjectID *int64
}

// SubscriptionDeleteParams controls subscription deletes.
type SubscriptionDeleteParams struct {
	ProjectID *int64
}
