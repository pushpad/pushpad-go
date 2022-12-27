package pushpad

import (
  "testing"
)

func TestConfigure(t *testing.T) {
  Configure("AUTH_TOKEN", "PROJECT_ID")

  if pushpadAuthToken != "AUTH_TOKEN" {
    t.Errorf("got %q instead of AUTH_TOKEN", pushpadAuthToken)
  }
  
  if pushpadProjectID != "PROJECT_ID" {
    t.Errorf("got %q instead of PROJECT_ID", pushpadProjectID)
  }
}
