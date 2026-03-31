package goforth

import (
	"errors"
)

var ErrEmptyRoom = errors.New("room must not be empty")

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

func (w *World) AddRoom(r *Room) error {
	if r.ID == "" {
		return ErrEmptyRoom
	}
	w.rooms[r.ID] = r
	return nil
}
