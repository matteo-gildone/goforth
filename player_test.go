package goforth

import "testing"

func TestNewPlayer(t *testing.T) {
	p := NewPlayer("entrance")

	if p.CurrentRoom() != "entrance" {
		t.Errorf("want: %q, got: %q", "entrance", p.CurrentRoom())
	}
}

func TestPlayer_MoveTo(t *testing.T) {
	p := NewPlayer("entrance")

	p.MoveTo("attic")

	if p.CurrentRoom() != "attic" {
		t.Errorf("want: %q, got: %q", "attic", p.CurrentRoom())
	}
}
