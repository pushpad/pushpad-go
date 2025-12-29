package pushpad

import (
	"testing"
)

func TestSignatureFor(t *testing.T) {
	Configure("5374d7dfeffa2eb49965624ba7596a09", 123)

	got := SignatureFor("user12345")
	want := "6627820dab00a1971f2a6d3ff16a5ad8ba4048a02b2d402820afc61aefd0b69f"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
