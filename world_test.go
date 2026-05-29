package goforth

import (
	"errors"
	"testing"
)

func TestWorld_AddRoom(t *testing.T) {
	entrance := NewRoom("entrance", "Entrance")

	w := NewWorld()
	err := w.AddRooms(entrance)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.RoomByID(entrance.ID)
	if !ok {
		t.Errorf("expected to find room %q", entrance.ID)
	}
}

func TestWorld_AddRooms(t *testing.T) {
	rooms := []*Room{
		NewRoom("entrance", "Entrance"),
		NewRoom("dining", "Dining room"),
	}

	w := NewWorld()
	err := w.AddRooms(rooms...)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	for _, v := range rooms {
		_, ok := w.RoomByID(v.ID)
		if !ok {
			t.Errorf("expected to find room %q", v.ID)
		}
	}
}

func TestWorld_AddRoomError(t *testing.T) {
	w := NewWorld()
	err := w.AddRooms(&Room{})
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidRoom) {
		t.Errorf("expected %v, got %v", ErrInvalidRoom, err)
	}
}

func TestWorld_AddRoomsError(t *testing.T) {
	w := NewWorld()
	err := w.AddRooms(NewRoom("entrance", "Entrance"), &Room{}, NewRoom("dining", "Dining room"))
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidRoom) {
		t.Errorf("expected %v, got %v", ErrInvalidRoom, err)
	}
	_, ok := w.RoomByID("entrance")
	if ok {
		t.Errorf("expected to not find room %q", "entrance")
	}

	_, ok = w.RoomByID("dining")
	if ok {
		t.Errorf("expected to not find %q", "dining")
	}
}

func TestWorld_AddObject(t *testing.T) {
	sword := NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true)

	w := NewWorld()
	err := w.AddObjects(sword)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	_, ok := w.ObjectByID(sword.ID)
	if !ok {
		t.Errorf("expected to find room %q", sword.ID)
	}
}

func TestWorld_AddObjectError(t *testing.T) {
	w := NewWorld()

	err := w.AddObjects(&Object{})
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidObject) {
		t.Errorf("expected %v, got %v", ErrInvalidObject, err)
	}
}

func TestWorld_AddObjects(t *testing.T) {
	objects := []*Object{
		NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
		NewObject("shield", "Shield", "A cracked shield bearing the crest of a fallen house", true),
	}

	w := NewWorld()
	err := w.AddObjects(objects...)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	for _, v := range objects {
		_, ok := w.ObjectByID(v.ID)
		if !ok {
			t.Errorf("expected to find object %q", v.ID)
		}
	}
}

func TestWorld_AddObjectsError(t *testing.T) {
	w := NewWorld()
	err := w.AddObjects(NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true), &Object{}, NewObject("shield", "Shield", "A cracked shield bearing the crest of a fallen house", true))
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidObject) {
		t.Errorf("expected %v, got %v", ErrInvalidObject, err)
	}
	_, ok := w.ObjectByID("sword")
	if ok {
		t.Errorf("expected to not find room %q", "sword")
	}

	_, ok = w.ObjectByID("shield")
	if ok {
		t.Errorf("expected to not find %q", "shield")
	}
}

func TestWorld_RoomByID_NotExistingID(t *testing.T) {
	w := NewWorld()
	_, ok := w.RoomByID("randomRoom")
	if ok {
		t.Errorf("expected room to not exist")
	}
}

func TestWorld_ObjectByID_NotExistingID(t *testing.T) {
	w := NewWorld()
	_, ok := w.ObjectByID("randomObject")
	if ok {
		t.Errorf("expected object to not exist")
	}
}

func TestWorld_PlaceObject(t *testing.T) {
	rooms := []*Room{
		NewRoom("entrance", "Entrance"),
	}
	objects := []*Object{
		NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
	}
	objectID := "sword"
	roomID := "entrance"
	w, err := setupWorld(rooms, objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = w.PlaceObject(objectID, roomID)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	objectsInRoom := w.ObjectsInRoom(roomID)

	if len(objectsInRoom) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(objectsInRoom))
	}
}

func TestWorld_PlaceObject_Errors(t *testing.T) {
	tests := []struct {
		name          string
		objectID      string
		roomID        string
		rooms         []*Room
		objects       []*Object
		wantErrorType string
		wantErrID     string
	}{
		{
			name:     "place not existing object in a room",
			objectID: "shield",
			roomID:   "entrance",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
			},
			wantErrorType: "object",
			wantErrID:     "shield",
		},
		{
			name:     "place object in a not existing room",
			objectID: "sword",
			roomID:   "sport",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
			},
			wantErrorType: "room",
			wantErrID:     "sport",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := setupWorld(tt.rooms, tt.objects)
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
				if roomErr.ID != tt.wantErrID {
					t.Errorf("want: %v, got: %v", tt.wantErrID, roomErr.ID)
				}
			case "object":
				if !errors.As(err, &objErr) {
					t.Fatalf("expected ObjectNotFoundErr, got %T", err)
				}
				if objErr.ID != tt.wantErrID {
					t.Errorf("want: %v, got: %v", tt.wantErrID, objErr.ID)
				}
			}

		})
	}
}

func TestWorld_ObjectsInRoom(t *testing.T) {
	tests := []struct {
		name          string
		roomID        string
		rooms         []*Room
		objects       []*Object
		objectsInRoom map[string]string
		wantLength    int
	}{
		{
			name:   "room has no objects",
			roomID: "entrance",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
			},
			objectsInRoom: map[string]string{},
			wantLength:    0,
		},
		{
			name:   "room has one object",
			roomID: "entrance",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
			},
			objectsInRoom: map[string]string{
				"sword": "entrance",
			},
			wantLength: 1,
		},
		{
			name:   "room has two objects",
			roomID: "entrance",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
				NewObject("shield", "Shield", "A cracked shield bearing the crest of a fallen house", true),
			},
			objectsInRoom: map[string]string{
				"sword":  "entrance",
				"shield": "entrance",
			},
			wantLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := setupWorld(tt.rooms, tt.objects)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			for k, v := range tt.objectsInRoom {
				err := w.PlaceObject(k, v)
				if err != nil {
					t.Fatalf("expected no error got: %v", err)
				}
			}

			objects := w.ObjectsInRoom(tt.roomID)

			if len(objects) != tt.wantLength {
				t.Errorf("want: %d, got: %d", tt.wantLength, len(objects))
			}
		})
	}
}

func TestWorld_PlayerInventory(t *testing.T) {
	tests := []struct {
		name               string
		rooms              []*Room
		objects            []*Object
		objectsInRoom      map[string]string
		objectsInInventory []string
		wantLength         int
	}{
		{
			name: "player has no objects",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
			},
			objectsInRoom:      map[string]string{},
			objectsInInventory: []string{},
			wantLength:         0,
		},
		{
			name: "player has one object",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
				NewObject("shield", "Shield", "A cracked shield bearing the crest of a fallen house", true),
			},
			objectsInRoom: map[string]string{
				"sword": "entrance",
			},
			objectsInInventory: []string{
				"shield",
			},
			wantLength: 1,
		},
		{
			name: "player has two objects",
			rooms: []*Room{
				NewRoom("entrance", "Entrance"),
			},
			objects: []*Object{
				NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
				NewObject("shield", "Shield", "A cracked shield bearing the crest of a fallen house", true),
				NewObject("key", "Dwarven key", "A heavy iron key, cold to the touch and stained with old blood", true),
				NewObject("potion", "Health potion", "A vial of crimson liquid that pulses faintly in the dark", true),
				NewObject("mana", "Mana potion", "A swirling blue draught that hums with arcane energy", true),
			},
			objectsInRoom: map[string]string{
				"sword":  "entrance",
				"shield": "entrance",
				"potion": "entrance",
			},
			objectsInInventory: []string{
				"key",
				"mana",
			},
			wantLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := setupWorld(tt.rooms, tt.objects)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			for k, v := range tt.objectsInRoom {
				err := w.PlaceObject(k, v)
				if err != nil {
					t.Fatalf("expected no error got: %v", err)
				}
			}

			for _, o := range tt.objectsInInventory {
				err := w.MoveObjectToPlayer(o)
				if err != nil {
					t.Fatalf("expected no error got: %v", err)
				}
			}

			objects := w.PlayerInventory()

			if len(objects) != tt.wantLength {
				t.Errorf("want: %d, got: %d", tt.wantLength, len(objects))
			}
		})
	}
}

func TestWorld_MoveObjectToPlayer(t *testing.T) {
	rooms := []*Room{
		NewRoom("entrance", "Entrance"),
	}
	objects := []*Object{
		NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
	}

	w, err := setupWorld(rooms, objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = w.MoveObjectToPlayer("sword")
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if !w.PlayerHasObject("sword") {
		t.Errorf("expected player to own %q", "sword")
	}
}

func TestWorld_MoveObjectToPlayer_Error(t *testing.T) {
	rooms := []*Room{
		NewRoom("entrance", "Entrance"),
	}
	objects := []*Object{
		NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true),
	}

	w, err := setupWorld(rooms, objects)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = w.MoveObjectToPlayer("shield")
	if err == nil {
		t.Fatal("expected error got nil")
	}

	var objErr *ObjectNotFoundErr
	if !errors.As(err, &objErr) {
		t.Fatalf("expected ObjectNotFoundErr, got %T", err)
	}
	if objErr.ID != "shield" {
		t.Errorf("want: %v, got: %v", "shield", objErr.ID)
	}
}

func setupWorld(rooms []*Room, objects []*Object) (*World, error) {
	w := NewWorld()
	err := w.AddRooms(rooms...)
	if err != nil {
		return nil, err
	}

	err = w.AddObjects(objects...)
	if err != nil {
		return nil, err
	}
	return w, nil
}
