package goforth

import (
	"errors"
	"testing"
)

func TestWorld_AddRoom(t *testing.T) {
	roomId := "entrance"
	roomDescription := "Main entrance of the castle"

	w, err := setupWorldWithRooms(map[string]string{roomId: roomDescription})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if len(w.rooms) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(w.rooms))
	}
}

func TestWorld_AddRoom_MultipleRooms(t *testing.T) {
	rooms := map[string]string{
		"entrance": "Entrance",
		"dining":   "Dining room",
	}

	w, err := setupWorldWithRooms(rooms)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
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
	if !errors.Is(err, ErrInvalidRoom) {
		t.Errorf("expected %v, got %v", ErrInvalidRoom, err)
	}
}

func TestWorld_AddObject(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	w, err := setupWorldWithObjects(map[string]string{objectId: objectName})
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if len(w.objects) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(w.objects))
	}
}

func TestWorld_AddRoom_MultipleObjects(t *testing.T) {
	objects := map[string]string{
		"sword":  "Sword",
		"shield": "Shield",
	}

	w, err := setupWorldWithObjects(objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if len(w.objects) != 2 {
		t.Errorf("want: %d, got: %d", 2, len(w.objects))
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
		t.Errorf("expected room to exists")
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
		t.Errorf("expected room to exists")
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
		wantErr   string
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
			wantErr: "failed to find current room ID \"library\"",
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
			wantErr: "failed to find adjacent room ID \"library\"",
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

			if err.Error() != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err.Error())
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
			return &World{}, err
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
			return &World{}, err
		}
	}
	return w, nil
}
