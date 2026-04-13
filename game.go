package goforth

import (
	"bufio"
	"errors"
	"io"
)

// Game holds the complete state of a running game session.
// World, Player and Registry are exported so consumers can inspect
// and extend state between commands if needed.
type Game struct {
	World    *World
	Player   *Player
	Registry *CommandRegistry
}

// Run reads lines from r in a loop, calls Parse, calls Registry.Dispatch.
// Stops on ErrQuit or an unexpected error.
func (g *Game) Run(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		cmd := Parse(line)
		if err := g.Registry.Dispatch(cmd, g); err != nil {
			if errors.Is(err, ErrQuit) {
				break
			}
			return err
		}
	}
	return nil
}

// NewGame creates a new game with the given world, player and registry.
func NewGame(w *World, p *Player, r *CommandRegistry) *Game {
	return &Game{
		World:    w,
		Player:   p,
		Registry: r,
	}
}
