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

	err := w.AddRoom(goforth.NewRoom("entrance", "Main entrance"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.AddRoom(goforth.NewRoom("dining", "Dining room"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	err = w.AddRoom(goforth.NewRoom("sport", "Sport room"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	err = w.AddRoom(goforth.NewRoom("library", "Library room"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.AddObject(goforth.NewObject("sword", "An elven sword"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.AddObject(goforth.NewObject("key", "A magic key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.ConnectRoomsBidirectional("entrance", goforth.North, "dining")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.ConnectRoomsBidirectional("dinging", goforth.West, "sport")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.ConnectRoomsBidirectional("dinging", goforth.North, "library")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.PlaceObject("sword", "library")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = w.PlaceObject("key", "sport")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	g := goforth.NewGame(w, p, r)
	err = g.Run(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
