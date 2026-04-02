package goforth

import (
	"errors"
	"testing"
)

func TestWorld_AddRoom_MultipleRooms(t *testing.T) {
	rooms := map[string]string{
		"entrance": "Entrance",
		"dining":   "Dining room",
	}

	w, err := setupWorldWithRooms(rooms)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	for k, _ := range rooms {
		_, ok := w.RoomByID(k)
		if !ok {
			t.Errorf("expected to find room %q", k)
		}
	}

}

func TestWorld_AddRoomError(t *testing.T) {
	w := NewWorld()

	err := w.AddRoom(&Room{})
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidRoom) {
		t.Errorf("expected %v, got %v", ErrInvalidRoom, err)
	}
}

func TestWorld_AddObject_MultipleObjects(t *testing.T) {
	objects := map[string]string{
		"sword":  "Sword",
		"shield": "Shield",
	}

	w, err := setupWorldWithObjects(objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	for k, _ := range objects {
		_, ok := w.ObjectByID(k)
		if !ok {
			t.Errorf("expected to find object %q", k)
		}
	}
}

func TestWorld_AddObjectError(t *testing.T) {
	w := NewWorld()

	err := w.AddObject(&Object{})
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidObject) {
		t.Errorf("expected %v, got %v", ErrInvalidObject, err)
	}
}

func TestWorld_RoomByID(t *testing.T) {
	roomId := "entrance"
	roomDescription := "Main entrance of the castle"

	w, err := setupWorldWithRooms(map[string]string{roomId: roomDescription})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	room, ok := w.RoomByID(roomId)
	if !ok {
		t.Errorf("expected room %q to exists", roomId)
	}

	if room.ID != roomId {
		t.Errorf("want: %q, got: %q", "entrance", room.ID)
	}

	if room.Description != roomDescription {
		t.Errorf("want: %q, got: %q", roomDescription, room.Description)
	}
}

func TestWorld_RoomByID_NotExistingID(t *testing.T) {
	roomId := "entrance"
	roomDescription := "Main entrance of the castle"

	w, err := setupWorldWithRooms(map[string]string{roomId: roomDescription})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.RoomByID("randomRoom")
	if ok {
		t.Errorf("expected room to not exists")
	}
}

func TestWorld_ObjectByID(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	w, err := setupWorldWithObjects(map[string]string{objectId: objectName})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	object, ok := w.ObjectByID(objectId)
	if !ok {
		t.Errorf("expected object %q to exists", objectId)
	}

	if object.ID != objectId {
		t.Errorf("want: %q, got: %q", "entrance", object.ID)
	}

	if object.Name != objectName {
		t.Errorf("want: %q, got: %q", objectName, object.Name)
	}
}

func TestWorld_ObjectByID_NotExistingID(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	w, err := setupWorldWithObjects(map[string]string{objectId: objectName})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.ObjectByID("randomObject")
	if ok {
		t.Errorf("expected room to not exists")
	}
}

func TestWorld_ConnectRooms(t *testing.T) {
	tests := []struct {
		name      string
		fromID    string
		toID      string
		direction Direction
		rooms     map[string]string
	}{
		{
			name:      "connecting entrance to dining room",
			fromID:    "entrance",
			toID:      "dining",
			direction: East,
			rooms: map[string]string{
				"entrance": "Entrance",
				"dining":   "Dining room",
			},
		},
		{
			name:      "connecting entrance to sport room",
			fromID:    "entrance",
			toID:      "sport",
			direction: West,
			rooms: map[string]string{
				"entrance": "Entrance",
				"sport":    "Sport room",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := setupWorldWithRooms(tt.rooms)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			err = w.ConnectRooms(tt.fromID, tt.direction, tt.toID)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			room, ok := w.RoomByID(tt.fromID)
			if !ok {
				t.Fatalf("expected room to exists")
			}

			if room.Exits[tt.direction] != tt.toID {
				t.Errorf("want: %q, got: %q", tt.toID, room.Exits[tt.direction])
			}
		})
	}
}

func TestWorld_ConnectRooms_Errors(t *testing.T) {
	tests := []struct {
		name      string
		fromID    string
		toID      string
		direction Direction
		rooms     map[string]string
		wantErrID string
	}{
		{
			name:      "connecting entrance to dining room",
			fromID:    "library",
			toID:      "dining",
			direction: East,
			rooms: map[string]string{
				"entrance": "Entrance",
				"dining":   "Dining room",
			},
			wantErrID: "library",
		},
		{
			name:      "connecting entrance to sport room",
			fromID:    "entrance",
			toID:      "library",
			direction: West,
			rooms: map[string]string{
				"entrance": "Entrance",
				"sport":    "Sport room",
			},
			wantErrID: "library",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := setupWorldWithRooms(tt.rooms)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			err = w.ConnectRooms(tt.fromID, tt.direction, tt.toID)
			if err == nil {
				t.Fatal("expected error got nil")
			}

			var roomErr *RoomNotFoundErr
			if !errors.As(err, &roomErr) {
				t.Fatalf("expected RoomNotFoundErr, got %T", err)
			}

			if roomErr.ID != tt.wantErrID {
				t.Errorf("want: %v, got: %v", tt.wantErrID, roomErr.ID)
			}
		})
	}
}

func setupWorldWithRooms(rooms map[string]string) (*World, error) {
	w := NewWorld()
	for k, v := range rooms {
		r := NewRoom(k, v)

		err := w.AddRoom(r)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

func setupWorldWithObjects(objects map[string]string) (*World, error) {
	w := NewWorld()
	for k, v := range objects {
		r := NewObject(k, v)

		err := w.AddObject(r)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}
