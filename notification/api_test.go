package notification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/pushpad/pushpad-go"
)

func TestListNotifications(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects/123/notifications").
		MatchParam("page", "2").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`[{"id":1,"body":"Hi"}]`)

	pushpad.Configure("TOKEN", 0)
	notifications, err := List(&NotificationListParams{ProjectID: pushpad.Int64(123), Page: pushpad.Int64(2)})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifications))
	}
	if notifications[0].ID != 1 {
		t.Errorf("expected notification ID 1, got %d", notifications[0].ID)
	}
}

func TestListNotificationsDefaultPage(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects/123/notifications").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`[]`)

	pushpad.Configure("TOKEN", 0)
	notifications, err := List(&NotificationListParams{ProjectID: pushpad.Int64(123)})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if len(notifications) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(notifications))
	}
}

func TestCreateNotification(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/notifications").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(201).
		BodyString(`{"id":99,"scheduled":10}`)

	pushpad.Configure("TOKEN", 0)
	response, err := Create(&NotificationCreateParams{ProjectID: pushpad.Int64(123), Body: pushpad.String("Hello")})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if response.ID != 99 {
		t.Errorf("expected notification ID 99, got %d", response.ID)
	}
	if response.Scheduled == nil || *response.Scheduled != 10 {
		t.Errorf("expected scheduled count 10, got %v", response.Scheduled)
	}
}

func TestCreateNotificationWithAllFields(t *testing.T) {
	defer gock.Off()

	sendAt, err := time.Parse(time.RFC3339Nano, "2016-07-06T10:09:00.000Z")
	if err != nil {
		t.Fatalf("expected no error parsing send_at, got %s", err)
	}

	params := NotificationCreateParams{
		ProjectID:          pushpad.Int64(123),
		Title:              pushpad.String("Foo Bar"),
		Body:               pushpad.String("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
		TargetURL:          pushpad.String("https://example.com"),
		IconURL:            pushpad.String("https://example.com/assets/icon.png"),
		BadgeURL:           pushpad.String("https://example.com/assets/badge.png"),
		ImageURL:           pushpad.String("https://example.com/assets/image.png"),
		TTL:                pushpad.Int64(604800),
		RequireInteraction: pushpad.Bool(false),
		Silent:             pushpad.Bool(false),
		Urgent:             pushpad.Bool(false),
		CustomData:         pushpad.String(""),
		Actions: Actions(
			NotificationAction{
				Title:     pushpad.String("A button"),
				TargetURL: pushpad.String("https://example.com/button-link"),
				Icon:      pushpad.String("https://example.com/assets/button-icon.png"),
				Action:    pushpad.String("myActionName"),
			},
		),
		Starred:       pushpad.Bool(false),
		SendAt:        pushpad.Time(sendAt),
		CustomMetrics: pushpad.Strings([]string{"metric1", "metric2"}),
		UIDs:          pushpad.Strings([]string{"uid0", "uid1", "uidN"}),
		Tags:          pushpad.Strings([]string{"tag1", "tagA && !tagB"}),
	}

	notificationJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/notifications").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		BodyString(string(notificationJSON)).
		Reply(201).
		BodyString(`{"id":123456789,"scheduled":9876}`)

	pushpad.Configure("TOKEN", 0)
	response, err := Create(&params)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if response.ID != 123456789 {
		t.Errorf("expected notification ID 123456789, got %d", response.ID)
	}
	if response.Scheduled == nil || *response.Scheduled != 9876 {
		t.Errorf("expected scheduled count 9876, got %v", response.Scheduled)
	}
}

func TestCreateNotificationMissingBody(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/notifications").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(422).
		BodyString(`{"error":"validation error"}`)

	pushpad.Configure("TOKEN", 0)
	_, err := Create(&NotificationCreateParams{ProjectID: pushpad.Int64(123)})
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

func TestNotificationSend(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects/123/notifications").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Accept", "application/json").
		MatchHeader("Authorization", "Bearer AUTH_TOKEN").
		Reply(201).
		BodyString("{\"id\": 123456789, \"scheduled\": 98765}")

	pushpad.Configure("AUTH_TOKEN", 0)

	n := NotificationCreateParams{ProjectID: pushpad.Int64(123), Body: pushpad.String("Hello world!")}
	res, err := Send(&n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	if res.ID != 123456789 {
		t.Errorf("got ID: %d, want ID: 123456789", res.ID)
	}
}

func TestNotificationWithUIDs(t *testing.T) {
	n := NotificationCreateParams{Body: pushpad.String("Hello user1"), UIDs: pushpad.Strings([]string{"user1"})}
	notificationJSON, err := json.Marshal(n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	got := string(notificationJSON)
	want := `{"body":"Hello user1","uids":["user1"],"tags":null}`

	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestNotificationWithTags(t *testing.T) {
	n := NotificationCreateParams{Body: pushpad.String("Hello tag1"), Tags: pushpad.Strings([]string{"tag1"})}
	notificationJSON, err := json.Marshal(n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	got := string(notificationJSON)
	want := `{"body":"Hello tag1","uids":null,"tags":["tag1"]}`

	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestGetNotification(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/notifications/77").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":77,"body":"Hello"}`)

	pushpad.Configure("TOKEN", 123)
	notification, err := Get(77, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if notification.ID != 77 {
		t.Errorf("expected notification ID 77, got %d", notification.ID)
	}
}

func TestCancelNotification(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Delete("/api/v1/notifications/77/cancel").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(204)

	pushpad.Configure("TOKEN", 123)
	if err := Cancel(77, nil); err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}

func TestListNotificationsMissingProjectID(t *testing.T) {
	pushpad.Configure("TOKEN", 0)
	_, err := List(nil)
	if err == nil || err.Error() != "pushpad: project ID is required" {
		t.Fatalf("expected project ID required error, got %v", err)
	}
}
