package main

import (
	"fmt"
	"os"

	"github.com/matteo-gildone/goforth"
)

func main() {
	entrance := goforth.NewRoom("entrance", "A gloomy entrance hall")
	dining := goforth.NewRoom("dining", "A dining room with a long oak table")
	sport := goforth.NewRoom("sport", "A sport room with trophies on the walls")
	library := goforth.NewRoom("library", "A dusty library lined with old books")
	entrance.North(dining).East(sport).West(library)
	dining.South(entrance)
	sport.West(entrance)
	library.East(entrance)
	sword := goforth.NewObject("sword", "An elven sword")
	key := goforth.NewObject("key", "A magic key")
	w := goforth.NewWorld()
	must(w.AddRooms(entrance, dining, sport, library))
	must(w.AddObjects(sword, key))
	must(w.PlaceObject("sword", "library"))
	must(w.PlaceObject("key", "sport"))
	mustValid(w.Validate())
	p := goforth.NewPlayer("entrance")
	r := goforth.NewCommandRegistry()
	goforth.RegisterDefaultHandlers(r)
	g := goforth.NewGame(w, p, r, os.Stdout)
	must(g.Run(os.Stdin))
}
func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func mustValid(issues []goforth.ValidationIssue) {
	for _, issue := range issues {
		fmt.Fprintf(os.Stderr, "%s\n", issue.Message)
	}
	if len(issues) > 0 {
		os.Exit(1)
	}
}
