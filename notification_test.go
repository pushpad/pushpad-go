package pushpad

import (
  "testing"
  "github.com/h2non/gock"
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
