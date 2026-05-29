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
	exits       map[Direction]Exit
	OnEnter     func(g *Game)
}

// North sets a bidirectional exit from r to to in the north direction.
// Returns r to allow chaining multiple exits.
func (r *Room) North(to *Room) *Room {
	r.exits[DirectionNorth] = Exit{RoomID: to.ID}
	to.exits[DirectionSouth] = Exit{RoomID: r.ID}
	return r
}

// South sets a bidirectional exit from r to to in the south direction.
// Returns r to allow chaining multiple exits.
func (r *Room) South(to *Room) *Room {
	r.exits[DirectionSouth] = Exit{RoomID: to.ID}
	to.exits[DirectionNorth] = Exit{RoomID: r.ID}
	return r
}

// East sets a bidirectional exit from r to to in the east direction.
// Returns r to allow chaining multiple exits.
func (r *Room) East(to *Room) *Room {
	r.exits[DirectionEast] = Exit{RoomID: to.ID}
	to.exits[DirectionWest] = Exit{RoomID: r.ID}
	return r
}

// West sets a bidirectional exit from r to to in the west direction.
// Returns r to allow chaining multiple exits.
func (r *Room) West(to *Room) *Room {
	r.exits[DirectionWest] = Exit{RoomID: to.ID}
	to.exits[DirectionEast] = Exit{RoomID: r.ID}
	return r
}

// Up sets a bidirectional exit from r to to in the up direction.
// Returns r to allow chaining multiple exits.
func (r *Room) Up(to *Room) *Room {
	r.exits[DirectionUp] = Exit{RoomID: to.ID}
	to.exits[DirectionDown] = Exit{RoomID: r.ID}
	return r
}

// Down sets a bidirectional exit from r to to in the down direction.
// Returns r to allow chaining multiple exits.
func (r *Room) Down(to *Room) *Room {
	r.exits[DirectionDown] = Exit{RoomID: to.ID}
	to.exits[DirectionUp] = Exit{RoomID: r.ID}
	return r
}

func (r *Room) Exit(dir Direction) (string, bool) {
	exit, ok := r.exits[dir]
	return exit.RoomID, ok
}

func (r *Room) ExitDirections() []Direction {
	dirs := make([]Direction, 0, len(r.exits))
	for dir := range r.exits {
		dirs = append(dirs, dir)
	}
	return dirs
}

// NewRoom creates a new room.
func NewRoom(id, description string) *Room {
	return &Room{
		ID:          id,
		Description: description,
		exits:       make(map[Direction]Exit),
	}
}
