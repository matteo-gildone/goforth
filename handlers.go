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
	if len(args) == 0 {
		fmt.Println("go where?")
		return nil
	}

	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	roomID, ok := currentRoom.Exits[Direction(args[0])]
	if !ok {
		fmt.Printf("you can't go %s\n", Direction(args[0]))
		return nil
	}

	g.Player.MoveTo(roomID)

	return nil
}

// TakeHandler moves named object from current room to player inventory.
// Register it as "take" in a CommandRegistry.
func TakeHandler(args []string, g *Game) error {
	if len(args) == 0 {
		fmt.Println("take what?")
		return nil
	}

	_, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	err := g.World.MoveObjectToPlayer(args[0])
	if err != nil {
		return err
	}

	return nil
}

// DropHandler moves named object from player inventory to current room.
// Register it as "drop" in a CommandRegistry.
func DropHandler(args []string, g *Game) error {
	if len(args) == 0 {
		fmt.Println("take what?")
		return nil
	}

	_, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	err := g.World.PlaceObject(args[0], g.Player.CurrentRoom())
	if err != nil {
		return err
	}

	return nil
}
