package goforth

// Direction represents a single direction in the game world.
type Direction string

const (
	DirectionNorth Direction = "north"
	DirectionSouth Direction = "south"
	DirectionEast  Direction = "east"
	DirectionWest  Direction = "west"
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
)

var oppositeDirectionMap = map[Direction]Direction{
	DirectionNorth: DirectionSouth,
	DirectionSouth: DirectionNorth,
	DirectionEast:  DirectionWest,
	DirectionWest:  DirectionEast,
	DirectionUp:    DirectionDown,
	DirectionDown:  DirectionUp,
}
