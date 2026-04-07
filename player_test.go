package goforth

import "testing"

func TestNewPlayer(t *testing.T) {
	p := NewPlayer("entrance")

	if p.CurrentRoom() != "entrance" {
		t.Errorf("want: %q, got: %q", "entrance", p.CurrentRoom())
	}
}
