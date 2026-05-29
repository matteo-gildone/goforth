package goforth

import (
	"errors"
	"fmt"
)

// ErrQuit is returned by a handler to signal a clean game exit.
// The game loop treats ErrQuit as a normal terminal condition, not an error.
var ErrQuit = errors.New("quit")

// Command represents a parsed player instruction.
// Name is the command verb. Args are the remaining tokens.
type Command struct {
	Name string
	Args []string
}

// HandlerFunc is the signature for all command handlers.
// args contains the tokens following the command verb.
// Returning ErrQuit signals the game loop to exit cleanly.
// Use this type when implementing custom commands:
//
//	r.Register("shout", goforth.HandlerFunc(func(args []string, g *goforth.Game) error {
//	    fmt.Fprintln(g.Out, "AAAARGH!")
//	    return nil
//	}))
type HandlerFunc func(args []string, g *Game) error

// CommandRegistry maps command names to their handlers.
// Consumers register their own commands alongside the built-in ones.
type CommandRegistry struct {
	handlers map[string]HandlerFunc
}

// Register adds a handler for the given command name.
// If name is empty the call is a no-op.
func (cr *CommandRegistry) Register(name string, h HandlerFunc) {
	if name != "" {
		cr.handlers[name] = h
	}
}

// Dispatch runs handler function associated to a command.
func (cr *CommandRegistry) Dispatch(cmd Command, g *Game) error {
	if cmd.Name == "" {
		return nil
	}

	command, ok := cr.handlers[cmd.Name]
	if !ok {
		fmt.Fprintln(g.Out, "I don't know how to do that.")
		return nil
	}

	return command(cmd.Args, g)
}

// NewCommandRegistry creates a new command registry.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]HandlerFunc),
	}
}
