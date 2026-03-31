package goforth

import (
	"errors"
	"testing"
)

func TestWorld_AddRoom(t *testing.T) {
	w := NewWorld()
	roomId := "entrance"
	roomDescription := "Main entrance of the castle"
	r := NewRoom(roomId, roomDescription)

	err := w.AddRoom(r)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if len(w.rooms) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(w.rooms))
	}

	if w.rooms["entrance"].ID != roomId {
		t.Errorf("want: %q, got: %q", "entrance", w.rooms["entrance"].ID)
	}

	if w.rooms["entrance"].Description != roomDescription {
		t.Errorf("want: %q, got: %q", roomDescription, w.rooms["entrance"].Description)
	}
}

func TestWorld_AddRooms(t *testing.T) {
	w := NewWorld()

	rooms := map[string]string{
		"entrance": "Entrance",
		"dining":   "Dining room",
	}

	for k, v := range rooms {
		r := NewRoom(k, v)

		err := w.AddRoom(r)
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}
	}

	if len(w.rooms) != 2 {
		t.Errorf("want: %d, got: %d", 2, len(w.rooms))
	}

}

func TestWorld_AddRoomError(t *testing.T) {
	w := NewWorld()

	err := w.AddRoom(&Room{})
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrEmptyRoom) {
		t.Errorf("expected %v, got %v", ErrEmptyRoom, err)
	}
}
