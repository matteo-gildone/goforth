package goforth

import (
	"errors"
	"fmt"
)

var ErrInvalidRoom = errors.New("room ID must not be empty")
var ErrInvalidObject = errors.New("object ID must not be empty")

// World represents the game world.
type World struct {
	rooms   map[string]*Room
	objects map[string]*Object
}

// NewWorld creates a new game world
func NewWorld() *World {
	return &World{
		rooms:   make(map[string]*Room),
		objects: make(map[string]*Object),
	}
}

// AddRoom adds a new room to the game world
func (w *World) AddRoom(r *Room) error {
	if r.ID == "" {
		return ErrInvalidRoom
	}
	w.rooms[r.ID] = r
	return nil
}

// AddObject adds a new object to the game world
func (w *World) AddObject(o *Object) error {
	if o.ID == "" {
		return ErrInvalidObject
	}
	w.objects[o.ID] = o
	return nil
}

// RoomByID gets room by id
func (w *World) RoomByID(id string) (*Room, bool) {
	room, ok := w.rooms[id]
	return room, ok
}

// ObjectByID gets object by id
func (w *World) ObjectByID(id string) (*Object, bool) {
	object, ok := w.objects[id]
	return object, ok
}

// ConnectRooms connects rooms together
func (w *World) ConnectRooms(fromID string, dir Direction, toID string) error {
	currentRoom, ok := w.rooms[fromID]
	if !ok {
		return fmt.Errorf("failed to find current room ID %q", fromID)
	}

	_, ok = w.rooms[toID]
	if !ok {
		return fmt.Errorf("failed to find adjacent room ID %q", toID)
	}

	currentRoom.Exits[dir] = toID
	return nil
}
