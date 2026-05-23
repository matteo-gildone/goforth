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

// ValidationKind represents the type of world validation issue found.
type ValidationKind string

const (
	// ValidationKindNoExits is returned when a room has no exits.
	ValidationKindNoExits ValidationKind = "no_exits"
	// ValidationKindUnregisteredExit is returned when a room exit points to a room not registered in the world.
	ValidationKindUnregisteredExit ValidationKind = "unregistered_room"
	// ValidationKindOneWay is returned when a room exit has no reverse exit.
	ValidationKindOneWay ValidationKind = "one_way"
)

// ValidationIssue describe different type of validation issue found by Validate
type ValidationIssue struct {
	RoomID  string
	Kind    ValidationKind
	Message string
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

// AddRoom adds r to the game world, making it available by ID.
// It returns ErrInvalidRoom if r has an empty ID.
func (w *World) AddRoom(r *Room) error {
	if r.ID == "" {
		return ErrInvalidRoom
	}
	w.rooms[r.ID] = r
	return nil
}

// AddRooms adds multiple rooms to the game world, making it available by ID.
// It returns ErrInvalidRoom if r has an empty ID.
func (w *World) AddRooms(rooms ...*Room) error {
	for _, r := range rooms {
		if r.ID == "" {
			return ErrInvalidRoom
		}
		w.rooms[r.ID] = r
	}

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

// AddObjects adds multiple objects to the game world.
// It returns ErrInvalidObject if o has an empty ID.
func (w *World) AddObjects(objects ...*Object) error {
	for _, o := range objects {
		if o.ID == "" {
			return ErrInvalidObject
		}
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

// Validate checks the world for configuration issues.
// Returns a slice of ValidationIssue describing each problem found.
// An empty slice means the world is valid.
func (w *World) Validate() []ValidationIssue {
	var issues []ValidationIssue
	for _, room := range w.rooms {
		if len(room.Exits) == 0 {
			issues = append(issues, ValidationIssue{room.ID, ValidationKindNoExits, fmt.Sprintf("room %q has no exits", room.ID)})
		}

		for direction, exit := range room.Exits {
			exitRoom, ok := w.RoomByID(exit.RoomID)
			oppositeDirection := oppositeDirectionMap[direction]
			if !ok {
				issues = append(issues, ValidationIssue{room.ID, ValidationKindUnregisteredExit, fmt.Sprintf("room %q does not exist", exit.RoomID)})
				continue
			}

			if room.ID != exitRoom.Exits[oppositeDirection].RoomID {
				issues = append(issues, ValidationIssue{room.ID, ValidationKindOneWay, fmt.Sprintf("room %q exit %q to %q has no reverse exit", room.ID, direction, exit.RoomID)})
			}
		}
	}
	return issues
}
