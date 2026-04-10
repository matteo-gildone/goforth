package goforth

import (
	"errors"
	"testing"
)

func TestHandlers_Error_InvalidRoom(t *testing.T) {
	tests := []struct {
		name string
		fn   HandlerFunc
		args []string
	}{
		{
			name: "look handler",
			fn:   LookHandler,
			args: []string{},
		},
		{
			name: "go handler",
			fn:   GoHandler,
			args: []string{"north"},
		},
		{
			name: "take handler",
			fn:   TakeHandler,
			args: []string{"sword"},
		},
		{
			name: "drop handler",
			fn:   DropHandler,
			args: []string{"sword"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := setupGame()
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			g.Player.MoveTo("wrong-room")

			err = tt.fn(tt.args, g)
			if err == nil {
				t.Fatal("expected error got nil")
			}

			var roomErr *RoomNotFoundErr
			if !errors.As(err, &roomErr) {
				t.Errorf("expected RoomNotFoundErr, got %T", err)
			}
			if roomErr.ID != "wrong-room" {
				t.Errorf("want: %v, got: %v", "wrong-room", roomErr.ID)
			}
		})

	}
}

func TestGoHandler(t *testing.T) {
	g, err := setupGame()
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = GoHandler([]string{"north"}, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if g.Player.CurrentRoom() != "dining" {
		t.Errorf("want: %q, got: %q", "dining", g.Player.CurrentRoom())
	}
}

func TestTakeHandler(t *testing.T) {
	g, err := setupGame()
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = TakeHandler([]string{"sword"}, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if !g.World.PlayerHasObject("sword") {
		t.Errorf("player should have %q", "sword")
	}
}

func TestDropHandler(t *testing.T) {
	g, err := setupGame()
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = g.World.MoveObjectToPlayer("sword")
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = DropHandler([]string{"sword"}, g)

	if g.World.PlayerHasObject("sword") {
		t.Errorf("player should have dropped %q", "sword")
	}
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
	c := NewCommandRegistry()

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

	w.ConnectRoomsBidirectional("entrance", North, "dining")

	g := NewGame(w, p, c)

	return g, nil
}
