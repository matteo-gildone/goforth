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
	w := NewWorld()
	w, err := setupWorldWithRooms(w, rooms)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	for k := range rooms {
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
	w := NewWorld()
	w, err := setupWorldWithObjects(w, objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	for k := range objects {
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
	w := NewWorld()
	w, err := setupWorldWithRooms(w, map[string]string{roomId: roomDescription})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.RoomByID(roomId)
	if !ok {
		t.Errorf("expected room %q to exist", roomId)
	}
}

func TestWorld_RoomByID_NotExistingID(t *testing.T) {
	roomId := "entrance"
	roomDescription := "Main entrance of the castle"
	w := NewWorld()
	w, err := setupWorldWithRooms(w, map[string]string{roomId: roomDescription})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.RoomByID("randomRoom")
	if ok {
		t.Errorf("expected room to not exist")
	}
}

func TestWorld_ObjectByID(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	w := NewWorld()
	w, err := setupWorldWithObjects(w, map[string]string{objectId: objectName})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.ObjectByID(objectId)
	if !ok {
		t.Errorf("expected object %q to exist", objectId)
	}
}

func TestWorld_ObjectByID_NotExistingID(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	w := NewWorld()
	w, err := setupWorldWithObjects(w, map[string]string{objectId: objectName})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.ObjectByID("randomObject")
	if ok {
		t.Errorf("expected object to not exist")
	}
}

func TestWorld_ConnectRooms(t *testing.T) {
	fromID := "entrance"
	toID := "dining"
	direction := East
	rooms := map[string]string{
		"entrance": "Entrance",
		"dining":   "Dining room",
	}
	w := NewWorld()
	w, err := setupWorldWithRooms(w, rooms)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = w.ConnectRooms(fromID, direction, toID)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	room, ok := w.RoomByID(fromID)
	if !ok {
		t.Fatalf("expected room to exist")
	}

	if room.Exits[direction] != toID {
		t.Errorf("want: %q, got: %q", toID, room.Exits[direction])
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
			w := NewWorld()
			w, err := setupWorldWithRooms(w, tt.rooms)
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

func TestWorld_PlaceObject(t *testing.T) {
	rooms := map[string]string{
		"entrance": "Entrance",
	}
	objects := map[string]string{
		"sword": "Sword",
	}
	objectID := "sword"
	roomID := "entrance"
	w := NewWorld()
	w, err := setupWorldWithRooms(w, rooms)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	w, err = setupWorldWithObjects(w, objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = w.PlaceObject(objectID, roomID)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.objectLocations[objectID]

	if !ok {
		t.Errorf("object: %q, expected to be in room: %q", objectID, roomID)
	}
}

func TestWorld_PlaceObject_Errors(t *testing.T) {
	tests := []struct {
		name          string
		objectID      string
		roomID        string
		direction     Direction
		rooms         map[string]string
		objects       map[string]string
		wantErrorType string
	}{
		{
			name:      "place object in a room",
			objectID:  "shield",
			roomID:    "entrance",
			direction: East,
			rooms: map[string]string{
				"entrance": "Entrance",
			},
			objects: map[string]string{
				"sword": "Sword",
			},
			wantErrorType: "object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()
			w, err := setupWorldWithRooms(w, tt.rooms)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			w, err = setupWorldWithObjects(w, tt.objects)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			err = w.PlaceObject(tt.objectID, tt.roomID)
			if err == nil {
				t.Fatal("expected error got nil")
			}

			var roomErr *RoomNotFoundErr
			var objErr *ObjectNotFoundErr

			switch tt.wantErrorType {
			case "room":
				if !errors.As(err, &roomErr) {
					t.Fatalf("expected RoomNotFoundErr, got %T", err)
				}
			case "object":
				if !errors.As(err, &objErr) {
					t.Fatalf("expected ObjectNotFoundErr, got %T", err)
				}
			}

		})
	}
}

func setupWorldWithRooms(w *World, rooms map[string]string) (*World, error) {
	for k, v := range rooms {
		r := NewRoom(k, v)

		err := w.AddRoom(r)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

func setupWorldWithObjects(w *World, objects map[string]string) (*World, error) {
	for k, v := range objects {
		r := NewObject(k, v)

		err := w.AddObject(r)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}
