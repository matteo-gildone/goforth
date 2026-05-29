package goforth

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidRoom   = errors.New("room ID must not be empty")
	ErrInvalidObject = errors.New("object ID must not be empty")
)

// RoomNotFoundErr is returned when a room ID cannot be resolved in the world.
type RoomNotFoundErr struct {
	ID string
}

func (e *RoomNotFoundErr) Error() string {
	return fmt.Sprintf("room %q not found", e.ID)
}

// ObjectNotFoundErr is returned when an object ID cannot be resolved in the world.
type ObjectNotFoundErr struct {
	ID string
}

func (e *ObjectNotFoundErr) Error() string {
	return fmt.Sprintf("object %q not found", e.ID)
}

// World represents the game world.
type World struct {
	rooms           map[string]*Room
	objects         map[string]*Object
	objectLocations map[string]string
}

// NewWorld creates a new game world.
func NewWorld() *World {
	return &World{
		rooms:           make(map[string]*Room),
		objects:         make(map[string]*Object),
		objectLocations: make(map[string]string),
	}
}

// AddRooms adds multiple rooms to the game world, making it available by ID.
// It returns ErrInvalidRoom if r has an empty ID.
func (w *World) AddRooms(rooms ...*Room) error {
	for _, r := range rooms {
		if r.ID == "" {
			return ErrInvalidRoom
		}
	}
	for _, r := range rooms {
		w.rooms[r.ID] = r
	}
	return nil
}

// AddObjects adds multiple objects to the game world.
// It returns ErrInvalidObject if o has an empty ID.
func (w *World) AddObjects(objects ...*Object) error {
	for _, o := range objects {
		if o.ID == "" {
			return ErrInvalidObject
		}
	}
	for _, o := range objects {
		w.objects[o.ID] = o
	}
	return nil
}

// RoomByID return the room with the given ID, or false if not found.
func (w *World) RoomByID(id string) (*Room, bool) {
	room, ok := w.rooms[id]
	return room, ok
}

// ObjectByID return the object with the given ID, or false if not found.
func (w *World) ObjectByID(id string) (*Object, bool) {
	object, ok := w.objects[id]
	return object, ok
}

// PlaceObject sets the initial location of an object within the world.
// It returns an error if the object or room ID is not recognized.
func (w *World) PlaceObject(objectID, roomID string) error {
	_, ok := w.ObjectByID(objectID)
	if !ok {
		return &ObjectNotFoundErr{ID: objectID}
	}
	_, ok = w.RoomByID(roomID)
	if !ok {
		return &RoomNotFoundErr{ID: roomID}
	}

	w.objectLocations[objectID] = roomID
	return nil
}

// ObjectsInRoom return the list of objects in a room
func (w *World) ObjectsInRoom(roomID string) []*Object {
	objects := make([]*Object, 0)
	for k, v := range w.objectLocations {
		if v == roomID {
			o, _ := w.ObjectByID(k)
			objects = append(objects, o)
		}
	}
	return objects
}

// PlayerInventory return the list of objects in player's inventory
func (w *World) PlayerInventory() []*Object {
	objects := make([]*Object, 0)
	for k := range w.objectLocations {
		if w.PlayerHasObject(k) {
			o, _ := w.ObjectByID(k)
			objects = append(objects, o)
		}
	}
	return objects
}

// MoveObjectToPlayer assign an object to a player
// It returns an error if the object ID is not recognized.
func (w *World) MoveObjectToPlayer(objectID string) error {
	_, ok := w.ObjectByID(objectID)
	if !ok {
		return &ObjectNotFoundErr{ID: objectID}
	}

	w.objectLocations[objectID] = "player"
	return nil
}

// PlayerHasObject check if player has an object with a certain ID
func (w *World) PlayerHasObject(objectID string) bool {
	return w.objectLocations[objectID] == "player"
}
