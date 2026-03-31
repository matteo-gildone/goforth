package goforth

type Direction string

const (
	North Direction = "north"
	South Direction = "south"
	East  Direction = "east"
	West  Direction = "west"
	Up    Direction = "up"
	Down  Direction = "down"
)

var aliasMap = map[string]Direction{
	"n": North,
	"s": South,
	"e": East,
	"w": West,
	"u": Up,
	"d": Down,
}

var oppositeDirectionMap = map[Direction]Direction{
	North: South,
	South: North,
	East:  West,
	West:  East,
	Up:    Down,
	Down:  Up,
}
