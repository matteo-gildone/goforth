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
// An empty or whitespace-only line returns a zero Command.
func Parse(line string) Command {
	line = strings.Trim(line, " ")
	line = strings.ToLower(line)
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return Command{}
	}
	return Command{Name: fields[0], Args: fields[1:]}
}
