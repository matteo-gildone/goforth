package goforth

import (
	"strings"
)

// Command represents a parsed player instruction.
// Name is the command verb. Args are the remaining tokens.
type Command struct {
	Name string
	Args []string
}

// Parse tokenizes a raw input line into a Command.
// Input is trimmed, lowercased, and split on whitespace.
// An empty or whitespace-only line returns a zero Command with nil Args.
func Parse(line string) Command {
	fields := strings.Fields(strings.ToLower(line))
	if len(fields) == 0 {
		return Command{}
	}
	return Command{Name: fields[0], Args: fields[1:]}
}
