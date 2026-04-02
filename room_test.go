package goforth

import (
	"testing"
)

func TestNewRoom(t *testing.T) {
	roomId := "entrance"
	roomDescription := "Main entrance of the castle"
	r := NewRoom(roomId, roomDescription)

	if r.ID != roomId {
		t.Errorf("want: %q, got: %q", "entrance", r.ID)
	}

	if r.Description != roomDescription {
		t.Errorf("want: %q, got: %q", roomDescription, r.Description)
	}
}
