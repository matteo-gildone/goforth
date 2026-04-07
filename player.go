package goforth

// Player represents the person playing the game.
// Player has no knowledge of the world — movement validation
// is the responsibility of the game loop.
type Player struct {
	currentRoomID string
}

func (p *Player) MoveTo(roomID string) {
	p.currentRoomID = roomID
}

func (p *Player) CurrentRoom() string {
	return p.currentRoomID
}

func NewPlayer(startingRoomID string) *Player {
	return &Player{currentRoomID: startingRoomID}
}
