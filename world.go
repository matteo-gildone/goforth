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

// World represents the game world.
type World struct {
	rooms   map[string]*Room
	objects map[string]*Object
}

// NewWorld creates a new game world.
func NewWorld() *World {
	return &World{
		rooms:   make(map[string]*Room),
		objects: make(map[string]*Object),
	}
}

// AddRoom adds r to the game world, making it available by ID.
// It returns ErrInvalidRoom if r has an empty ID.
func (w *World) AddRoom(r *Room) error {
	if r.ID == "" {
		return ErrInvalidRoom
	}
	w.rooms[r.ID] = r
	return nil
}

// AddObject adds o to the game world.
// It returns ErrInvalidObject if o has an empty ID.
func (w *World) AddObject(o *Object) error {
	if o.ID == "" {
		return ErrInvalidObject
	}
	w.objects[o.ID] = o
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

// ConnectRooms adds a directional exit from one room to another.
// Returns RoomNotFoundErr if either rooms ID is not registered in the game world.
func (w *World) ConnectRooms(fromID string, dir Direction, toID string) error {
	currentRoom, ok := w.RoomByID(fromID)
	if !ok {
		return &RoomNotFoundErr{ID: fromID}
	}

	_, ok = w.RoomByID(toID)
	if !ok {
		return &RoomNotFoundErr{ID: toID}
	}

	currentRoom.Exits[dir] = toID
	return nil
}
