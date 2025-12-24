package subscription

import "time"

// Subscription represents a Pushpad subscription.
type Subscription struct {
	ID          int        `json:"id,omitempty"`
	ProjectID   int        `json:"project_id,omitempty"`
	Endpoint    string     `json:"endpoint,omitempty"`
	P256DH      string     `json:"p256dh,omitempty"`
	Auth        string     `json:"auth,omitempty"`
	UID         string     `json:"uid,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	LastClickAt *time.Time `json:"last_click_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

// SubscriptionCreateParams is the payload to create a subscription.
type SubscriptionCreateParams struct {
	ProjectID *int      `json:"-"`
	Endpoint  *string   `json:"endpoint"`
	P256DH    *string   `json:"p256dh,omitempty"`
	Auth      *string   `json:"auth,omitempty"`
	UID       *string   `json:"uid,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
}

// SubscriptionUpdateParams is the payload to update a subscription.
type SubscriptionUpdateParams struct {
	ProjectID *int      `json:"-"`
	UID       *string   `json:"uid,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
}

// SubscriptionListParams controls subscription listing.
type SubscriptionListParams struct {
	ProjectID *int
	Page      *int
	PerPage   *int
	UIDs      *[]string
	Tags      *[]string
}

// SubscriptionGetParams controls subscription fetches.
type SubscriptionGetParams struct {
	ProjectID *int
}

// SubscriptionDeleteParams controls subscription deletes.
type SubscriptionDeleteParams struct {
	ProjectID *int
}
