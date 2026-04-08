package goforth

type Game struct {
	World    *World
	Player   *Player
	Registry *CommandRegistry
}

func NewGame() *Game {
	return &Game{
		World:    NewWorld(),
		Player:   NewPlayer(""),
		Registry: NewCommandRegistry(),
	}
}
