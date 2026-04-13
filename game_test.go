package goforth

import (
	"strings"
	"testing"
)

func TestGame_Run(t *testing.T) {
	t.Run("player moves to another room", func(t *testing.T) {
		g, err := setupGame()
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		err = g.Run(strings.NewReader("look\ngo north\nquit"))
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		if g.Player.CurrentRoom() != "dining" {
			t.Errorf("want: %q, got: %q", "dining", g.Player.CurrentRoom())
		}
	})

	t.Run("player pick up sword", func(t *testing.T) {
		g, err := setupGame()
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		err = g.Run(strings.NewReader("look\ntake sword\nquit"))
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		if !g.World.PlayerHasObject("sword") {
			t.Errorf("player should have %q", "sword")
		}
	})
	t.Run("player drop item in another room", func(t *testing.T) {
		g, err := setupGame()
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		err = g.Run(strings.NewReader("look\ntake sword\ngo north\ndrop sword\nquit"))
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}

		if g.World.PlayerHasObject("sword") {
			t.Errorf("player shouldn't have %q", "sword")
		}

		if len(g.World.ObjectsInRoom("dining")) != 1 {
			t.Errorf("want: %d, got: %d", 1, len(g.World.ObjectsInRoom("dining")))
		}

		if len(g.World.ObjectsInRoom("entrance")) != 0 {
			t.Errorf("want: %d, got: %d", 0, len(g.World.ObjectsInRoom("entrance")))
		}
	})
}

func setupGame() (*Game, error) {
	rooms := map[string]string{
		"entrance": "Entrance",
		"dining":   "Dining room",
	}
	objects := map[string]string{
		"sword":  "Sword",
		"shield": "Shield",
		"key":    "Dwarven key",
		"potion": "Health potion",
		"mana":   "Mana potion",
	}

	p := NewPlayer("entrance")
	w := NewWorld()
	r := NewCommandRegistry()

	RegisterDefaultHandlers(r)

	for k, v := range rooms {
		r := NewRoom(k, v)

		err := w.AddRoom(r)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range objects {
		r := NewObject(k, v)

		err := w.AddObject(r)
		if err != nil {
			return nil, err
		}
	}

	err := w.ConnectRoomsBidirectional("entrance", North, "dining")
	if err != nil {
		return nil, err
	}

	g := NewGame(w, p, r)
	err = g.World.PlaceObject("sword", "entrance")
	if err != nil {
		return nil, err
	}

	return g, nil
}
