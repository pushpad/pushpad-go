package project

import (
	"encoding/json"
	"testing"

	"github.com/h2non/gock"
	"github.com/pushpad/pushpad-go"
)

func TestListProjects(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`[{"id":1,"name":"Main"}]`)

	pushpad.Configure("TOKEN", 123)
	projects, err := List(nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].ID != 1 {
		t.Errorf("expected project ID 1, got %d", projects[0].ID)
	}
}

func TestCreateProject(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(201).
		BodyString(`{"id":2,"name":"New Project","website":"https://example.com","sender_id":9}`)

	pushpad.Configure("TOKEN", 123)
	payload := &ProjectCreateParams{
		SenderID: pushpad.Int64(9),
		Name:     pushpad.String("New Project"),
		Website:  pushpad.String("https://example.com"),
	}
	project, err := Create(payload)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if project.ID != 2 {
		t.Errorf("expected project ID 2, got %d", project.ID)
	}
}

func TestCreateProjectWithAllFields(t *testing.T) {
	defer gock.Off()

	params := ProjectCreateParams{
		SenderID:                     pushpad.Int64(98765),
		Name:                         pushpad.String("My Project"),
		Website:                      pushpad.String("https://example.com"),
		IconURL:                      pushpad.String("https://example.com/icon.png"),
		BadgeURL:                     pushpad.String("https://example.com/badge.png"),
		NotificationsTTL:             pushpad.Int64(604800),
		NotificationsRequireInteract: pushpad.Bool(false),
		NotificationsSilent:          pushpad.Bool(false),
	}

	projectJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		BodyString(string(projectJSON)).
		Reply(201).
		BodyString(`{"id":12345,"name":"My Project","website":"https://example.com","sender_id":98765}`)

	pushpad.Configure("TOKEN", 123)
	project, err := Create(&params)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if project.ID != 12345 {
		t.Errorf("expected project ID 12345, got %d", project.ID)
	}
}

func TestCreateProjectMissingFields(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(422).
		BodyString(`{"error":"validation error"}`)

	pushpad.Configure("TOKEN", 123)
	_, err := Create(&ProjectCreateParams{})
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

func TestGetProject(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Get("/api/v1/projects/2").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":2,"name":"New Project"}`)

	pushpad.Configure("TOKEN", 123)
	project, err := Get(2, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if project.ID != 2 {
		t.Errorf("expected project ID 2, got %d", project.ID)
	}
}

func TestUpdateProject(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Patch("/api/v1/projects/2").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(200).
		BodyString(`{"id":2,"name":"Updated Project"}`)

	pushpad.Configure("TOKEN", 123)
	update := &ProjectUpdateParams{Name: pushpad.String("Updated Project")}
	project, err := Update(2, update)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if project.Name != "Updated Project" {
		t.Errorf("expected updated name, got %q", project.Name)
	}
}

func TestDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Delete("/api/v1/projects/2").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(202)

	pushpad.Configure("TOKEN", 123)
	if err := Delete(2, nil); err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}

func TestAPIErrorOnServerFailure(t *testing.T) {
	defer gock.Off()

	gock.New("https://pushpad.xyz").
		Post("/api/v1/projects").
		MatchHeader("Authorization", "Bearer TOKEN").
		Reply(500).
		BodyString(`{"error":"boom"}`)

	pushpad.Configure("TOKEN", 123)
	_, err := Create(&ProjectCreateParams{
		SenderID: pushpad.Int64(1),
		Name:     pushpad.String("Failing Project"),
		Website:  pushpad.String("https://example.com"),
	})
	apiErr, ok := err.(*pushpad.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
}
