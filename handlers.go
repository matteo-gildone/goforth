package goforth

import (
	"fmt"
	"slices"
)

// LookHandler prints the description and exits of the player's current room.
// Register it as "look" in a CommandRegistry.
func LookHandler(args []string, g *Game) error {
	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	fmt.Fprintln(g.Out, currentRoom.Description)
	fmt.Fprintln(g.Out, "Exits:")
	for _, dir := range currentRoom.ExitDirections() {
		fmt.Fprintf(g.Out, "  %s\n", dir)
	}
	return nil
}

// GoHandler resolves direction and moves player.
// Register it as "go" in a CommandRegistry.
func GoHandler(args []string, g *Game) error {
	if len(args) == 0 {
		fmt.Fprintln(g.Out, "go where?")
		return nil
	}

	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	roomID, ok := currentRoom.Exit(Direction(args[0]))
	if !ok {
		fmt.Fprintf(g.Out, "you can't go %s\n", Direction(args[0]))
		return nil
	}

	g.Player.MoveTo(roomID)

	currentRoom, ok = g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	if currentRoom.OnEnter != nil {
		currentRoom.OnEnter(g)
	}
	return nil
}

// TakeHandler moves named object from current room to player inventory.
// Register it as "take" in a CommandRegistry.
func TakeHandler(args []string, g *Game) error {
	if len(args) == 0 {
		fmt.Fprintln(g.Out, "take what?")
		return nil
	}

	_, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	objects := g.World.ObjectsInRoom(g.Player.CurrentRoom())

	for _, object := range objects {
		if object.ID != args[0] {
			continue
		}

		// object exists in the room but cannot be picked up
		if !object.Takeable {
			fmt.Fprintln(g.Out, "you can't take that!")
			return nil
		}

		return g.World.MoveObjectToPlayer(args[0])
	}

	fmt.Fprintf(g.Out, "there is no %q here\n", args[0])
	return nil
}

// DropHandler moves named object from player inventory to current room.
// Register it as "drop" in a CommandRegistry.
func DropHandler(args []string, g *Game) error {
	if len(args) == 0 {
		fmt.Fprintln(g.Out, "drop what?")
		return nil
	}

	_, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
	}

	owned := g.World.PlayerHasObject(args[0])
	if !owned {
		fmt.Fprintf(g.Out, "you don't own %q\n", args[0])
		return nil
	}

	return g.World.PlaceObject(args[0], g.Player.CurrentRoom())
}

// InventoryHandler lists objects present in players' inventory.
// Register it as "inventory" in a CommandRegistry.
func InventoryHandler(args []string, g *Game) error {
	objects := g.World.PlayerInventory()
	if len(objects) == 0 {
		fmt.Fprintln(g.Out, "nothing in the inventory")
		return nil
	}

	fmt.Fprintln(g.Out, "Inventory:")
	for _, object := range objects {
		fmt.Fprintf(g.Out, "  %s\n", object.Name)
	}

	return nil
}

// QuitHandler signal end of the game.
// Register it as "quit" in a CommandRegistry.
func QuitHandler(args []string, g *Game) error {
	return ErrQuit
}

func ExamineHandler(args []string, g *Game) error {
	if len(args) == 0 {
		fmt.Fprintln(g.Out, "examine what?")
		return nil
	}

	inRoom := slices.ContainsFunc(g.World.ObjectsInRoom(g.Player.CurrentRoom()), func(object *Object) bool {
		return object.ID == args[0]
	})

	inInventory := slices.ContainsFunc(g.World.PlayerInventory(), func(object *Object) bool {
		return object.ID == args[0]
	})

	if !inRoom && !inInventory {
		fmt.Fprintln(g.Out, "you don't see that here")
		return nil
	}

	value, _ := g.World.ObjectByID(args[0])
	fmt.Fprintln(g.Out, value.Description)

	return nil
}

// RegisterDefaultHandlers register the default handlers.
func RegisterDefaultHandlers(r *CommandRegistry) {
	directionAliases := map[string]Direction{
		"n":     DirectionNorth,
		"north": DirectionNorth,
		"south": DirectionSouth,
		"s":     DirectionSouth,
		"west":  DirectionWest,
		"w":     DirectionWest,
		"east":  DirectionEast,
		"e":     DirectionEast,
		"up":    DirectionUp,
		"u":     DirectionUp,
		"down":  DirectionDown,
		"d":     DirectionDown,
	}

	for alias, dir := range directionAliases {
		d := dir
		r.Register(alias, func(args []string, g *Game) error {
			return GoHandler([]string{string(d)}, g)
		})
	}

	commandAliases := map[string]HandlerFunc{
		"i": InventoryHandler,
		"x": ExamineHandler,
	}

	for command, handler := range commandAliases {
		r.Register(command, handler)
	}

	r.Register("go", GoHandler)
	r.Register("look", LookHandler)
	r.Register("take", TakeHandler)
	r.Register("drop", DropHandler)
	r.Register("inventory", InventoryHandler)
	r.Register("examine", ExamineHandler)
	r.Register("quit", QuitHandler)
}
