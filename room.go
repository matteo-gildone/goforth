package goforth

// Room represents a single location in the game world.
// Room is connected to other rooms via directional exits.
type Room struct {
	ID          string
	Description string
	Exits       map[Direction]string
}

// NewRoom creates a new room.
func NewRoom(id, description string) *Room {
	return &Room{
		ID:          id,
		Description: description,
		Exits:       make(map[Direction]string),
	}
}
