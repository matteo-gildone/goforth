package goforth

import (
	"bytes"
	"strings"
	"testing"
)

func TestGame_Run(t *testing.T) {
	t.Run("player moves to another room", func(t *testing.T) {
		var buf bytes.Buffer
		g, err := setupGame(&buf)
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
		var buf bytes.Buffer
		g, err := setupGame(&buf)
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
		var buf bytes.Buffer
		g, err := setupGame(&buf)
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

func setupGame(out *bytes.Buffer) (*Game, error) {
	entrance := NewRoom("entrance", "Entrance")
	dining := NewRoom("dining", "Dining room")
	sword := NewObject("sword", "Sword", "A blade forged in shadow, its edge never dulls", true)
	shield := NewObject("shield", "Shield", "A cracked shield bearing the crest of a fallen house", true)
	key := NewObject("key", "Dwarven key", "A heavy iron key, cold to the touch and stained with old blood", true)
	potion := NewObject("potion", "Health potion", "A vial of crimson liquid that pulses faintly in the dark", true)
	mana := NewObject("mana", "Mana potion", "A swirling blue draught that hums with arcane energy", true)
	armor := NewObject("armor", "Armor", "A suit of blackened plate armour, too heavy to carry", false)

	p := NewPlayer("entrance")
	w := NewWorld()
	c := NewCommandRegistry()

	RegisterDefaultHandlers(c)

	err := w.AddObjects(sword, shield, key, potion, mana, armor)
	if err != nil {
		return nil, err
	}

	err = w.AddRooms(entrance, dining)
	if err != nil {
		return nil, err
	}

	entrance.North(dining)

	g := NewGame(w, p, c, out)
	err = g.World.PlaceObject("sword", "entrance")
	if err != nil {
		return nil, err
	}

	return g, nil
}

func setupAliasWorld(out *bytes.Buffer) (*Game, error) {
	entrance := NewRoom("entrance", "Entrance")
	dining := NewRoom("dining", "Dining room")
	sport := NewRoom("sport", "Sport room")
	library := NewRoom("library", "a dusty library")
	kitchen := NewRoom("kitchen", "a messy kitchen")

	p := NewPlayer("entrance")
	w := NewWorld()
	c := NewCommandRegistry()

	RegisterDefaultHandlers(c)

	err := w.AddRooms(entrance, dining, sport, library, kitchen)
	if err != nil {
		return nil, err
	}

	entrance.North(dining).East(sport).West(library).South(kitchen)

	g := NewGame(w, p, c, out)

	return g, nil
}
