package subscription

import (
	"encoding/json"
	"testing"

	"github.com/h2non/gock"
	"github.com/pushpad/pushpad-go"
)

func TestListSubscriptions(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects/123/subscriptions").
		MatchParam("page", "1").
		MatchParam("per_page", "20").
		MatchParam("uids[]", "u1").
		MatchParam("uids[]", "u2").
		MatchParam("tags[]", "tag1").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		SetHeader("X-Total-Count", "2").
		BodyString(`[{"id":10,"endpoint":"https://example.com/1"},{"id":11,"endpoint":"https://example.com/2"}]`)

	pushpad.Configure("TOKEN", 0)
	params := &SubscriptionListParams{
		ProjectID: pushpad.Int64(123),
		Page:      pushpad.Int64(1),
		PerPage:   pushpad.Int64(20),
		UIDs:      pushpad.Strings([]string{"u1", "u2"}),
		Tags:      pushpad.Strings([]string{"tag1"}),
	}
	subscriptions, total, err := List(params)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if total != 2 {
		t.Errorf("expected total count 2, got %d", total)
	}
	if len(subscriptions) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(subscriptions))
	}
	if subscriptions[0].ID != 10 {
		t.Errorf("expected subscription ID 10, got %d", subscriptions[0].ID)
	}
}

func TestListSubscriptionsNoOptions(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects/123/subscriptions").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`[]`)

	pushpad.Configure("TOKEN", 0)
	subscriptions, total, err := List(&SubscriptionListParams{ProjectID: pushpad.Int64(123)})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if total != 0 {
		t.Errorf("expected total count 0, got %d", total)
	}
	if len(subscriptions) != 0 {
		t.Fatalf("expected 0 subscriptions, got %d", len(subscriptions))
	}
}

func TestCreateSubscription(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/subscriptions").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(201).
		BodyString(`{"id":50,"endpoint":"https://example.com/1"}`)

	pushpad.Configure("TOKEN", 0)
	subscription, err := Create(&SubscriptionCreateParams{ProjectID: pushpad.Int64(123), Endpoint: pushpad.String("https://example.com/1")})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if subscription.ID != 50 {
		t.Errorf("expected subscription ID 50, got %d", subscription.ID)
	}
}

func TestCreateSubscriptionWithAllFields(t *testing.T) {
	defer gock.Off()

	params := SubscriptionCreateParams{
		ProjectID: pushpad.Int64(123),
		Endpoint:  pushpad.String("https://example.com/push/f7Q1Eyf7EyfAb1"),
		P256DH:    pushpad.String("BCQVDTlYWdl05lal3lG5SKr3VxTrEWpZErbkxWrzknHrIKFwihDoZpc_2sH6Sh08h-CacUYI-H8gW4jH-uMYZQ4="),
		Auth:      pushpad.String("cdKMlhgVeSPzCXZ3V7FtgQ=="),
		UID:       pushpad.String("user1"),
		Tags:      pushpad.Strings([]string{"tag1", "tag2"}),
	}

	subscriptionJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/subscriptions").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		BodyString(string(subscriptionJSON)).
		Reply(201).
		BodyString(`{"id":12345,"endpoint":"https://example.com/push/f7Q1Eyf7EyfAb1"}`)

	pushpad.Configure("TOKEN", 0)
	subscription, err := Create(&params)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if subscription.ID != 12345 {
		t.Errorf("expected subscription ID 12345, got %d", subscription.ID)
	}
}

func TestCreateSubscriptionMissingEndpoint(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/subscriptions").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(422).
		BodyString(`{"error":"validation error"}`)

	pushpad.Configure("TOKEN", 0)
	_, err := Create(&SubscriptionCreateParams{ProjectID: pushpad.Int64(123)})
	apiErr, ok := err.(*pushpad.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("expected status 422, got %d", apiErr.StatusCode)
	}
	if apiErr.Body != `{"error":"validation error"}` {
		t.Errorf("expected validation error body, got %q", apiErr.Body)
	}
}

func TestGetSubscription(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects/123/subscriptions/50").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":50,"endpoint":"https://example.com/1"}`)

	pushpad.Configure("TOKEN", 0)
	subscription, err := Get(50, &SubscriptionGetParams{ProjectID: pushpad.Int64(123)})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if subscription.ID != 50 {
		t.Errorf("expected subscription ID 50, got %d", subscription.ID)
	}
}

func TestUpdateSubscription(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Patch("/api/v1/projects/123/subscriptions/50").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":50,"uid":"new-user"}`)

	pushpad.Configure("TOKEN", 0)
	update := &SubscriptionUpdateParams{ProjectID: pushpad.Int64(123), UID: pushpad.String("new-user")}
	subscription, err := Update(50, update)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if subscription.UID != "new-user" {
		t.Errorf("expected uid new-user, got %q", subscription.UID)
	}
}

func TestDeleteSubscription(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Delete("/api/v1/projects/123/subscriptions/50").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(204)

	pushpad.Configure("TOKEN", 0)
	if err := Delete(50, &SubscriptionDeleteParams{ProjectID: pushpad.Int64(123)}); err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}
