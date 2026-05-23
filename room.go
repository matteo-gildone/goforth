package goforth

// Exit represents a single exit for a room
type Exit struct {
	RoomID string
}

// Room represents a single location in the game world.
// Room is connected to other rooms via directional exits.
type Room struct {
	ID          string
	Description string
	Exits       map[Direction]Exit
}

// North sets a one-directional exit from r to to in the north direction.
// Returns r to allow chaining multiple exits.
func (r *Room) North(to *Room) *Room {
	r.Exits[DirectionNorth] = Exit{RoomID: to.ID}
	return r
}

// South sets a one-directional exit from r to to in the south direction.
// Returns r to allow chaining multiple exits.
func (r *Room) South(to *Room) *Room {
	r.Exits[DirectionSouth] = Exit{RoomID: to.ID}
	return r
}

// East sets a one-directional exit from r to to in the east direction.
// Returns r to allow chaining multiple exits.
func (r *Room) East(to *Room) *Room {
	r.Exits[DirectionEast] = Exit{RoomID: to.ID}
	return r
}

// West sets a one-directional exit from r to to in the west direction.
// Returns r to allow chaining multiple exits.
func (r *Room) West(to *Room) *Room {
	r.Exits[DirectionWest] = Exit{RoomID: to.ID}
	return r
}

// Up sets a one-directional exit from r to to in the up direction.
// Returns r to allow chaining multiple exits.
func (r *Room) Up(to *Room) *Room {
	r.Exits[DirectionUp] = Exit{RoomID: to.ID}
	return r
}

// Down sets a one-directional exit from r to to in the down direction.
// Returns r to allow chaining multiple exits.
func (r *Room) Down(to *Room) *Room {
	r.Exits[DirectionDown] = Exit{RoomID: to.ID}
	return r
}

// NewRoom creates a new room.
func NewRoom(id, description string) *Room {
	return &Room{
		ID:          id,
		Description: description,
		Exits:       make(map[Direction]Exit),
	}
}
