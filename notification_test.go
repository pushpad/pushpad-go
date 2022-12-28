package pushpad

import (
  "testing"
  "github.com/h2non/gock"
  "encoding/json"
)

func TestNotificationSend(t *testing.T) {
  defer gock.Off()

  gock.New("https://pushpad.xyz").
    Post("/api/v1/projects/PROJECT_ID/notifications").
    MatchHeader("Content-Type", "application/json").
    MatchHeader("Accept", "application/json").
    MatchHeader("Authorization", "Token token=\"AUTH_TOKEN\"").
    Reply(201).
    BodyString("{\"id\": 123456789, \"scheduled\": 98765}")
  
  Configure("AUTH_TOKEN", "PROJECT_ID")
  
  n := Notification { Body: "Hello world!" }
  res, err := n.Send()

  if err != nil {
    t.Errorf("got an error: %s", err)
  }
  
  if res.ID != 123456789 {
    t.Errorf("got ID: %d, want ID: 123456789", res.ID)
  }
}

func TestNotificationWithUIDs(t *testing.T) {
  n := Notification { Body: "Hello user1", UIDs: []string{"user1"} }
  notificationJSON, err := json.Marshal(n)
  
  if err != nil {
    t.Errorf("got an error: %s", err)
  }
  
  got := string(notificationJSON)
  want := `{"body":"Hello user1","uids":["user1"],"tags":null}`

  if got != want {
    t.Errorf("got: %q, want: %q", got, want)
  }
}

func TestNotificationWithTags(t *testing.T) {
  n := Notification { Body: "Hello tag1", Tags: []string{"tag1"} }
  notificationJSON, err := json.Marshal(n)
  
  if err != nil {
    t.Errorf("got an error: %s", err)
  }
  
  got := string(notificationJSON)
  want := `{"body":"Hello tag1","uids":null,"tags":["tag1"]}`

  if got != want {
    t.Errorf("got: %q, want: %q", got, want)
  }
}
