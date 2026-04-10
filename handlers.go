package goforth

import "fmt"

// LookHandler prints the description and exits of the player's current room.
// Register it as "look" in a CommandRegistry.
func LookHandler(args []string, g *Game) error {
	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	fmt.Println(currentRoom.Description)
	fmt.Println("Exits:")
	for k := range currentRoom.Exits {
		fmt.Println(k)
	}
	return nil
}

// GoHandler resolves direction and moves player.
// Register it as "go" in a CommandRegistry.
func GoHandler(args []string, g *Game) error {
	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	roomID, ok := currentRoom.Exits[Direction(args[0])]
	if !ok {
		fmt.Printf("you can't go %s\n", Direction(args[0]))
	}

	_, ok = g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	g.Player.MoveTo(roomID)

	return nil
}
