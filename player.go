package goforth

// Player represents the person playing the game.
// Player has no knowledge of the world — movement validation
// is the responsibility of the game loop.
type Player struct {
	currentRoomID string
}

// MoveTo moves player into another room
func (p *Player) MoveTo(roomID string) {
	p.currentRoomID = roomID
}

// CurrentRoom return current player location
func (p *Player) CurrentRoom() string {
	return p.currentRoomID
}

// NewPlayer creates a new player in the game world.
func NewPlayer(startingRoomID string) *Player {
	return &Player{currentRoomID: startingRoomID}
}
