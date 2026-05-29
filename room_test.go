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

	if len(r.ExitDirections()) != 0 {
		t.Errorf("want: %d, got: %d", 0, len(r.ExitDirections()))
	}

	if r.OnEnter != nil {
		t.Error("expect no default initialisation")
	}
}
