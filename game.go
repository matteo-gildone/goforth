package goforth

type Game struct {
	World    *World
	Player   *Player
	Registry *CommandRegistry
}

func NewGame(w *World, p *Player, r *CommandRegistry) *Game {
	return &Game{
		World:    w,
		Player:   p,
		Registry: r,
	}
}
