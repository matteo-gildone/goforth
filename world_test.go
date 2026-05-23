package goforth

import (
	"errors"
	"testing"
)

func TestWorld_AddRoom(t *testing.T) {
	entrance := NewRoom("entrance", "Entrance")

	w := NewWorld()
	err := w.AddRoom(entrance)
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
	err := w.AddRoom(&Room{})
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
	if !ok {
		t.Errorf("expected to find room %q", "entrance")
	}

	_, ok = w.RoomByID("dining")
	if ok {
		t.Errorf("expected to not find %q", "dining")
	}
}

func TestWorld_AddObject(t *testing.T) {
	sword := NewObject("sword", "Sword")

	w := NewWorld()
	err := w.AddObject(sword)
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

	err := w.AddObject(&Object{})
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidObject) {
		t.Errorf("expected %v, got %v", ErrInvalidObject, err)
	}
}

func TestWorld_AddObjects(t *testing.T) {
	objects := []*Object{
		NewObject("sword", "Sword"),
		NewObject("shield", "Shield"),
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
	err := w.AddObjects(NewObject("sword", "Sword"), &Object{}, NewObject("shield", "Shield"))
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if !errors.Is(err, ErrInvalidObject) {
		t.Errorf("expected %v, got %v", ErrInvalidObject, err)
	}
	_, ok := w.ObjectByID("sword")
	if !ok {
		t.Errorf("expected to find room %q", "sword")
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
		NewObject("sword", "Sword"),
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
				NewObject("sword", "Sword"),
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
				NewObject("sword", "Sword"),
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
				NewObject("sword", "Sword"),
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
				NewObject("sword", "Sword"),
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
				NewObject("sword", "Sword"),
				NewObject("shield", "Shield"),
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
				NewObject("sword", "Sword"),
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
				NewObject("sword", "Sword"),
				NewObject("shield", "Shield"),
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
				NewObject("sword", "Sword"),
				NewObject("shield", "Shield"),
				NewObject("key", "Dwarven key"),
				NewObject("potion", "Health potion"),
				NewObject("mana", "Mana potion"),
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
		NewObject("sword", "Sword"),
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
		NewObject("sword", "Sword"),
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

func TestWorld_Validate(t *testing.T) {
	t.Run("room with no exit", func(t *testing.T) {
		entrance := NewRoom("entrance", "Entrance")

		w := NewWorld()
		err := w.AddRoom(entrance)
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}
		issues := w.Validate()
		if len(issues) != 1 {
			t.Errorf("want: %d, got: %d", 1, len(issues))
		}

		if issues[0].Kind != ValidationKindNoExits {
			t.Errorf("want: %q, got: %q", ValidationKindNoExits, issues[0].Kind)
		}

		if issues[0].RoomID != "entrance" {
			t.Errorf("want: %q, got: %q", "entrance", issues[0].RoomID)
		}
	})

	t.Run("room exit is not registered", func(t *testing.T) {
		entrance := NewRoom("entrance", "Entrance")

		w := NewWorld()
		err := w.AddRoom(entrance)
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		entrance.North(NewRoom("unregistered-room", "unregistered room"))

		issues := w.Validate()
		if len(issues) != 1 {
			t.Errorf("want: %d, got: %d", 1, len(issues))
		}

		if issues[0].Kind != ValidationKindUnregisteredExit {
			t.Errorf("want: %q, got: %q", ValidationKindUnregisteredExit, issues[0].Kind)
		}
	})

	t.Run("room is not bidirectional connected", func(t *testing.T) {
		entrance := NewRoom("entrance", "Entrance")
		dining := NewRoom("dining", "Dining room")
		unrelated := NewRoom("kitchen", "Kitchen")

		w := NewWorld()
		err := w.AddRooms(entrance, dining, unrelated)
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		dining.South(unrelated)
		unrelated.North(dining)
		entrance.North(dining)

		issues := w.Validate()

		if len(issues) > 2 {
			t.Errorf("want: %d, got: %d", 2, len(issues))
		}

		if issues[0].Kind != ValidationKindOneWay {
			t.Errorf("want: %q, got: %q", ValidationKindOneWay, issues[0].Kind)
		}
	})

	t.Run("valid world returns no issues", func(t *testing.T) {
		entrance := NewRoom("entrance", "Entrance")
		dining := NewRoom("dining", "Dining room")

		w := NewWorld()
		err := w.AddRooms(entrance, dining)
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		entrance.North(dining)
		dining.South(entrance)

		issues := w.Validate()

		if len(issues) != 0 {
			t.Errorf("want: %d, got: %d", 0, len(issues))
		}
	})
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
