package goforth

import "fmt"

// LookHandler prints the description and exits of the player's current room.
// Register it as "look" in a CommandRegistry.
func LookHandler(args []string, g *Game) error {
	fmt.Println("here!")
	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	fmt.Println(currentRoom.Description)
	fmt.Println("Exits:")
	for k := range currentRoom.Exits {
		fmt.Printf("  %s\n", k)
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
		fmt.Println("drop what?")
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

// InventoryHandler lists objects present in players' inventory.
// Register it as "inventory" in a CommandRegistry.
func InventoryHandler(args []string, g *Game) error {
	objects := g.World.PlayerInventory()
	if len(objects) == 0 {
		fmt.Println("nothing in the inventory")
		return nil
	}

	fmt.Println("Inventory:")
	for _, object := range objects {
		fmt.Printf("  %s\n", object.Name)
	}

	return nil
}

// QuitHandler signal end of the game.
// Register it as "quit" in a CommandRegistry.
func QuitHandler(args []string, g *Game) error {
	return ErrQuit
}

// RegisterDefaultHandlers register the default handlers.
func RegisterDefaultHandlers(r *CommandRegistry) {
	aliases := map[string]Direction{
		"n":     North,
		"north": North,
		"south": South,
		"s":     South,
		"west":  West,
		"w":     West,
		"east":  East,
		"e":     East,
		"up":    Up,
		"u":     Up,
		"Down":  Down,
		"d":     Down,
	}

	for alias, dir := range aliases {
		d := dir
		r.Register(alias, func(args []string, g *Game) error {
			return GoHandler([]string{string(d)}, g)
		})
	}

	r.Register("go", GoHandler)
	r.Register("look", LookHandler)
	r.Register("take", TakeHandler)
	r.Register("drop", DropHandler)
	r.Register("inventory", InventoryHandler)
	r.Register("quit", QuitHandler)
}
