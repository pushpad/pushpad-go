package pushpad

import (
	"testing"
)

func TestConfigure(t *testing.T) {
	Configure("AUTH_TOKEN", 123)

	if pushpadAuthToken != "AUTH_TOKEN" {
		t.Errorf("got %q instead of AUTH_TOKEN", pushpadAuthToken)
	}

	if pushpadProjectID != 123 {
		t.Errorf("got %d instead of project ID 123", pushpadProjectID)
	}
}
