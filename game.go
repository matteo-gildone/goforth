package goforth

// Game holds the complete state of a running game session
// World, Player and Registry are exported so consumers can inspect
// and extend state between commands if needed.
type Game struct {
	World    *World
	Player   *Player
	Registry *CommandRegistry
}

// NewGame creates a new game with the given world, player and registry.
func NewGame(w *World, p *Player, r *CommandRegistry) *Game {
	return &Game{
		World:    w,
		Player:   p,
		Registry: r,
	}
}
