package sender

import "time"

// Sender represents a Pushpad sender.
type Sender struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	VAPIDPrivateKey string     `json:"vapid_private_key"`
	VAPIDPublicKey  string     `json:"vapid_public_key"`
	CreatedAt       time.Time  `json:"created_at"`
}

// SenderCreateParams is the payload to create a sender.
type SenderCreateParams struct {
	Name            *string `json:"name,omitempty"`
	VAPIDPrivateKey *string `json:"vapid_private_key,omitempty"`
	VAPIDPublicKey  *string `json:"vapid_public_key,omitempty"`
}

// SenderUpdateParams is the payload to update a sender.
type SenderUpdateParams struct {
	Name *string `json:"name,omitempty"`
}

// SenderListParams controls sender listing.
type SenderListParams struct{}

// SenderGetParams controls sender fetches.
type SenderGetParams struct{}

// SenderDeleteParams controls sender deletes.
type SenderDeleteParams struct{}
