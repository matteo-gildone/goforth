package main

import (
	"fmt"
	"os"

	"github.com/matteo-gildone/goforth"
)

func main() {
	p := goforth.NewPlayer("entrance")
	w := goforth.NewWorld()
	r := goforth.NewCommandRegistry()

	goforth.RegisterDefaultHandlers(r)

	must(w.AddRoom(goforth.NewRoom("entrance", "Main entrance")))
	must(w.AddRoom(goforth.NewRoom("dining", "Dining room")))
	must(w.AddRoom(goforth.NewRoom("sport", "Sport room")))
	must(w.AddRoom(goforth.NewRoom("library", "Library room")))
	must(w.AddObject(goforth.NewObject("sword", "An elven sword")))
	must(w.AddObject(goforth.NewObject("key", "A magic key")))
	must(w.ConnectRoomsBidirectional("entrance", goforth.North, "dining"))
	must(w.ConnectRoomsBidirectional("dinging", goforth.West, "sport"))
	must(w.ConnectRoomsBidirectional("dinging", goforth.North, "library"))
	must(w.PlaceObject("sword", "library"))
	must(w.PlaceObject("key", "sport"))
	g := goforth.NewGame(w, p, r)
	must(g.Run(os.Stdin))
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
