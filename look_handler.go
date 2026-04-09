package goforth

import "fmt"

func LookHandler(args []string, g *Game) error {
	currentRoom, ok := g.World.RoomByID(g.Player.CurrentRoom())
	if !ok {
		return fmt.Errorf("%q doesn't exist", g.Player.CurrentRoom())
	}

	fmt.Println(currentRoom.Description)
	fmt.Println("Exits:")
	for k, _ := range currentRoom.Exits {
		fmt.Println(k)
	}
	return nil
}
