package notification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/pushpad/pushpad-go/v1"
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
	if response.Scheduled != 10 {
		t.Errorf("expected scheduled count 10, got %d", response.Scheduled)
	}
	if !response.SendAt.IsZero() {
		t.Errorf("expected send_at to be zero value, got %s", response.SendAt)
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
		Actions: &[]NotificationActionParams{
			{
				Title:     pushpad.String("A button"),
				TargetURL: pushpad.String("https://example.com/button-link"),
				Icon:      pushpad.String("https://example.com/assets/button-icon.png"),
				Action:    pushpad.String("myActionName"),
			},
		},
		Starred:       pushpad.Bool(false),
		SendAt:        pushpad.Time(sendAt),
		CustomMetrics: pushpad.StringSlice([]string{"metric1", "metric2"}),
		UIDs:          pushpad.StringSlice([]string{"uid0", "uid1", "uidN"}),
		Tags:          pushpad.StringSlice([]string{"tag1", "tagA && !tagB"}),
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
	if response.Scheduled != 9876 {
		t.Errorf("expected scheduled count 9876, got %d", response.Scheduled)
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
	n := NotificationCreateParams{Body: pushpad.String("Hello user1"), UIDs: pushpad.StringSlice([]string{"user1"})}
	notificationJSON, err := json.Marshal(n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	got := string(notificationJSON)
	want := `{"body":"Hello user1","uids":["user1"]}`

	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestNotificationWithEmptyUIDs(t *testing.T) {
	n := NotificationCreateParams{Body: pushpad.String("Hello user1"), UIDs: pushpad.StringSlice([]string{})}
	notificationJSON, err := json.Marshal(n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	got := string(notificationJSON)
	want := `{"body":"Hello user1","uids":[]}`

	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestNotificationWithTags(t *testing.T) {
	n := NotificationCreateParams{Body: pushpad.String("Hello tag1"), Tags: pushpad.StringSlice([]string{"tag1"})}
	notificationJSON, err := json.Marshal(n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	got := string(notificationJSON)
	want := `{"body":"Hello tag1","tags":["tag1"]}`

	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestNotificationWithEmptyTags(t *testing.T) {
	n := NotificationCreateParams{Body: pushpad.String("Hello tag1"), Tags: pushpad.StringSlice([]string{})}
	notificationJSON, err := json.Marshal(n)

	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	got := string(notificationJSON)
	want := `{"body":"Hello tag1","tags":[]}`

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

func TestGetNotificationWithAllFields(t *testing.T) {
	defer gock.Off()

	sendAt, err := time.Parse(time.RFC3339Nano, "2016-07-06T10:09:00.000Z")
	if err != nil {
		t.Fatalf("expected no error parsing send_at, got %s", err)
	}

	createdAt, err := time.Parse(time.RFC3339Nano, "2016-07-06T10:58:39.192Z")
	if err != nil {
		t.Fatalf("expected no error parsing created_at, got %s", err)
	}

	gock.New("https://pushpad.xyz").
		Get("/api/v1/notifications/123456789").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":123456789,"project_id":123,"title":"Foo Bar","body":"Lorem ipsum dolor sit amet, consectetur adipiscing elit.","target_url":"https://example.com","icon_url":"https://example.com/assets/icon.png","badge_url":"https://example.com/assets/badge.png","image_url":"https://example.com/assets/image.png","ttl":604800,"require_interaction":false,"silent":false,"urgent":false,"custom_data":"","actions":[{"title":"A button","target_url":"https://example.com/button-link","icon":"https://example.com/assets/button-icon.png","action":"myActionName"}],"starred":false,"send_at":"2016-07-06T10:09:00.000Z","custom_metrics":["metric1","metric2"],"uids":["uid0","uid1","uidN"],"tags":["tag1","tagA && !tagB"],"created_at":"2016-07-06T10:58:39.192Z","successfully_sent_count":4,"opened_count":1,"scheduled_count":400,"scheduled":true,"cancelled":false}`)

	pushpad.Configure("TOKEN", 123)
	notification, err := Get(123456789, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if notification.ID != 123456789 {
		t.Errorf("expected notification ID 123456789, got %d", notification.ID)
	}
	if notification.ProjectID != 123 {
		t.Errorf("expected project ID 123, got %d", notification.ProjectID)
	}
	if notification.Title != "Foo Bar" {
		t.Errorf("expected title Foo Bar, got %q", notification.Title)
	}
	if notification.Body != "Lorem ipsum dolor sit amet, consectetur adipiscing elit." {
		t.Errorf("expected body Lorem ipsum..., got %q", notification.Body)
	}
	if notification.TargetURL != "https://example.com" {
		t.Errorf("expected target_url https://example.com, got %q", notification.TargetURL)
	}
	if notification.IconURL != "https://example.com/assets/icon.png" {
		t.Errorf("expected icon_url https://example.com/assets/icon.png, got %q", notification.IconURL)
	}
	if notification.BadgeURL != "https://example.com/assets/badge.png" {
		t.Errorf("expected badge_url https://example.com/assets/badge.png, got %q", notification.BadgeURL)
	}
	if notification.ImageURL != "https://example.com/assets/image.png" {
		t.Errorf("expected image_url https://example.com/assets/image.png, got %q", notification.ImageURL)
	}
	if notification.TTL != 604800 {
		t.Errorf("expected ttl 604800, got %d", notification.TTL)
	}
	if notification.RequireInteraction != false {
		t.Errorf("expected require_interaction false, got %v", notification.RequireInteraction)
	}
	if notification.Silent != false {
		t.Errorf("expected silent false, got %v", notification.Silent)
	}
	if notification.Urgent != false {
		t.Errorf("expected urgent false, got %v", notification.Urgent)
	}
	if notification.CustomData != "" {
		t.Errorf("expected custom_data empty string, got %q", notification.CustomData)
	}
	if len(notification.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(notification.Actions))
	}
	if notification.Actions[0].Title != "A button" {
		t.Errorf("expected action title A button, got %v", notification.Actions[0].Title)
	}
	if notification.Actions[0].TargetURL != "https://example.com/button-link" {
		t.Errorf("expected action target_url https://example.com/button-link, got %v", notification.Actions[0].TargetURL)
	}
	if notification.Actions[0].Icon != "https://example.com/assets/button-icon.png" {
		t.Errorf("expected action icon https://example.com/assets/button-icon.png, got %v", notification.Actions[0].Icon)
	}
	if notification.Actions[0].Action != "myActionName" {
		t.Errorf("expected action myActionName, got %v", notification.Actions[0].Action)
	}
	if notification.Starred != false {
		t.Errorf("expected starred false, got %v", notification.Starred)
	}
	if !notification.SendAt.Equal(sendAt) {
		t.Errorf("expected send_at %s, got %s", sendAt.Format(time.RFC3339Nano), notification.SendAt.Format(time.RFC3339Nano))
	}
	if len(notification.CustomMetrics) != 2 || notification.CustomMetrics[0] != "metric1" || notification.CustomMetrics[1] != "metric2" {
		t.Errorf("expected custom_metrics [metric1 metric2], got %v", notification.CustomMetrics)
	}
	if len(notification.UIDs) != 3 || notification.UIDs[0] != "uid0" || notification.UIDs[1] != "uid1" || notification.UIDs[2] != "uidN" {
		t.Errorf("expected uids [uid0 uid1 uidN], got %v", notification.UIDs)
	}
	if len(notification.Tags) != 2 || notification.Tags[0] != "tag1" || notification.Tags[1] != "tagA && !tagB" {
		t.Errorf("expected tags [tag1 tagA && !tagB], got %v", notification.Tags)
	}
	if !notification.CreatedAt.Equal(createdAt) {
		t.Errorf("expected created_at %s, got %s", createdAt.Format(time.RFC3339Nano), notification.CreatedAt.Format(time.RFC3339Nano))
	}
	if notification.SuccessfullySent != 4 {
		t.Errorf("expected successfully_sent_count 4, got %d", notification.SuccessfullySent)
	}
	if notification.OpenedCount != 1 {
		t.Errorf("expected opened_count 1, got %d", notification.OpenedCount)
	}
	if notification.ScheduledCount != 400 {
		t.Errorf("expected scheduled_count 400, got %d", notification.ScheduledCount)
	}
	if notification.Scheduled != true {
		t.Errorf("expected scheduled true, got %v", notification.Scheduled)
	}
	if notification.Cancelled != false {
		t.Errorf("expected cancelled false, got %v", notification.Cancelled)
	}
}

func TestGetNotificationWithNullFields(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/notifications/88").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":88,"body":"Hello","image_url":null,"custom_data":null,"actions":null,"send_at":null,"custom_metrics":null,"uids":null,"tags":null,"scheduled_count":null,"scheduled":null}`)

	pushpad.Configure("TOKEN", 123)
	notification, err := Get(88, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if notification.ImageURL != "" {
		t.Errorf("expected image_url empty string, got %q", notification.ImageURL)
	}
	if notification.CustomData != "" {
		t.Errorf("expected custom_data empty string, got %q", notification.CustomData)
	}
	if notification.Actions != nil {
		t.Errorf("expected actions nil, got %v", notification.Actions)
	}
	if !notification.SendAt.IsZero() {
		t.Errorf("expected send_at zero value, got %s", notification.SendAt)
	}
	if notification.CustomMetrics != nil {
		t.Errorf("expected custom_metrics nil, got %v", notification.CustomMetrics)
	}
	if notification.UIDs != nil {
		t.Errorf("expected uids nil, got %v", notification.UIDs)
	}
	if notification.Tags != nil {
		t.Errorf("expected tags nil, got %v", notification.Tags)
	}
	if notification.ScheduledCount != 0 {
		t.Errorf("expected scheduled_count 0, got %d", notification.ScheduledCount)
	}
	if notification.Scheduled != false {
		t.Errorf("expected scheduled false, got %v", notification.Scheduled)
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
