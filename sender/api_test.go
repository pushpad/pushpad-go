package sender

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/pushpad/pushpad-go"
)

func TestListSenders(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/senders").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`[{"id":1,"name":"Sender"}]`)

	pushpad.Configure("TOKEN", 123)
	senders, err := List(nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if len(senders) != 1 {
		t.Fatalf("expected 1 sender, got %d", len(senders))
	}
	if senders[0].ID != 1 {
		t.Errorf("expected sender ID 1, got %d", senders[0].ID)
	}
}

func TestCreateSender(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/senders").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(201).
		BodyString(`{"id":5,"name":"New Sender"}`)

	pushpad.Configure("TOKEN", 123)
	sender, err := Create(&SenderCreateParams{Name: pushpad.String("New Sender")})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if sender.ID != 5 {
		t.Errorf("expected sender ID 5, got %d", sender.ID)
	}
}

func TestCreateSenderWithAllFields(t *testing.T) {
	defer gock.Off()

	params := SenderCreateParams{
		Name:            pushpad.String("My Sender"),
		VAPIDPrivateKey: pushpad.String("-----BEGIN EC PRIVATE KEY----- ..."),
		VAPIDPublicKey:  pushpad.String("-----BEGIN PUBLIC KEY----- ..."),
	}

	senderJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	gock.New("https://pushpad.xyz").
		Post("/api/v1/senders").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		BodyString(string(senderJSON)).
		Reply(201).
		BodyString(`{"id":12345,"name":"My Sender"}`)

	pushpad.Configure("TOKEN", 123)
	sender, err := Create(&params)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if sender.ID != 12345 {
		t.Errorf("expected sender ID 12345, got %d", sender.ID)
	}
}

func TestCreateSenderMissingName(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/senders").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(422).
		BodyString(`{"error":"validation error"}`)

	pushpad.Configure("TOKEN", 123)
	_, err := Create(&SenderCreateParams{})
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

func TestGetSender(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/senders/5").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":5,"name":"New Sender"}`)

	pushpad.Configure("TOKEN", 123)
	sender, err := Get(5, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if sender.ID != 5 {
		t.Errorf("expected sender ID 5, got %d", sender.ID)
	}
}

func TestGetSenderWithAllFields(t *testing.T) {
	defer gock.Off()

	createdAt, err := time.Parse(time.RFC3339Nano, "2016-07-06T10:58:39.192Z")
	if err != nil {
		t.Fatalf("expected no error parsing created_at, got %s", err)
	}

	gock.New("https://pushpad.xyz").
		Get("/api/v1/senders/98765").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":98765,"name":"My Sender","vapid_private_key":"-----BEGIN EC PRIVATE KEY----- ...","vapid_public_key":"-----BEGIN PUBLIC KEY----- ...","created_at":"2016-07-06T10:58:39.192Z"}`)

	pushpad.Configure("TOKEN", 123)
	sender, err := Get(98765, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if sender.ID != 98765 {
		t.Errorf("expected sender ID 98765, got %d", sender.ID)
	}
	if sender.Name != "My Sender" {
		t.Errorf("expected name My Sender, got %q", sender.Name)
	}
	if sender.VAPIDPrivateKey != "-----BEGIN EC PRIVATE KEY----- ..." {
		t.Errorf("expected VAPID private key value, got %q", sender.VAPIDPrivateKey)
	}
	if sender.VAPIDPublicKey != "-----BEGIN PUBLIC KEY----- ..." {
		t.Errorf("expected VAPID public key value, got %q", sender.VAPIDPublicKey)
	}
	if !sender.CreatedAt.Equal(createdAt) {
		t.Errorf("expected created_at %s, got %s", createdAt.Format(time.RFC3339Nano), sender.CreatedAt.Format(time.RFC3339Nano))
	}
}

func TestUpdateSender(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Patch("/api/v1/senders/5").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":5,"name":"Updated Sender"}`)

	pushpad.Configure("TOKEN", 123)
	update := &SenderUpdateParams{Name: pushpad.String("Updated Sender")}
	sender, err := Update(5, update)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if sender.Name != "Updated Sender" {
		t.Errorf("expected updated name, got %q", sender.Name)
	}
}

func TestDeleteSender(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Delete("/api/v1/senders/5").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(204)

	pushpad.Configure("TOKEN", 123)
	if err := Delete(5, nil); err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}
