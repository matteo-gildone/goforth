package goforth

import (
	"bytes"
	"errors"
	"slices"
	"strings"
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
			var buf bytes.Buffer
			g, err := setupGame(&buf)
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
	var buf bytes.Buffer
	g, err := setupGame(&buf)
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

func TestGoHandler_InvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{
			name:       "no args",
			args:       []string{},
			wantOutput: "go where?",
		},
		{
			name:       "invalid direction",
			args:       []string{"sideways"},
			wantOutput: "you can't go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			g, err := setupGame(&buf)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}
			err = GoHandler(tt.args, g)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}
			if !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("want output to contain %q, got %q", tt.wantOutput, buf.String())
			}
			if g.Player.CurrentRoom() != "entrance" {
				t.Errorf("player should not have moved, got %q", g.Player.CurrentRoom())
			}
		})
	}
}

func TestTakeHandler(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = g.World.PlaceObject("sword", "entrance")
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

func TestTakeHandler_InvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{
			name:       "no args",
			args:       []string{},
			wantOutput: "take what?",
		},
		{
			name:       "invalid object",
			args:       []string{"frying pan"},
			wantOutput: "there is no",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			g, err := setupGame(&buf)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}
			err = TakeHandler(tt.args, g)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}
			if !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("want output to contain %q, got %q", tt.wantOutput, buf.String())
			}
		})
	}
}

func TestDropHandler(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = g.World.MoveObjectToPlayer("sword")
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	err = DropHandler([]string{"sword"}, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	if g.World.PlayerHasObject("sword") {
		t.Errorf("player should have dropped %q", "sword")
	}

	if !slices.ContainsFunc(g.World.ObjectsInRoom("entrance"), func(object *Object) bool {
		return object.ID == "sword"
	}) {
		t.Errorf("%q should be in the room", "sword")
	}
}

func TestDropHandler_InvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{
			name:       "no args",
			args:       []string{},
			wantOutput: "drop what?",
		},
		{
			name:       "invalid direction",
			args:       []string{"frying pan"},
			wantOutput: "you don't own",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			g, err := setupGame(&buf)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}
			err = DropHandler(tt.args, g)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}
			if !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("want output to contain %q, got %q", tt.wantOutput, buf.String())
			}
		})
	}
}

func TestLookHandler(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	err = LookHandler([]string{}, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	if !strings.Contains(buf.String(), "Entrance") {
		t.Errorf("want output to contain %q, got %q", "Entrance", buf.String())
	}
	if !strings.Contains(buf.String(), "Exits") {
		t.Errorf("want output to contain %q, got %q", "Exits", buf.String())
	}
	if !strings.Contains(buf.String(), "north") {
		t.Errorf("want output to contain %q, got %q", "north", buf.String())
	}
}

func TestInventoryHandler_Empty(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	err = InventoryHandler([]string{}, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	if !strings.Contains(buf.String(), "nothing in the inventory") {
		t.Errorf("want output to contain %q, got %q", "nothing in the inventory", buf.String())
	}
}

func TestInventoryHandler_WithItems(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	err = g.World.MoveObjectToPlayer("sword")
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	err = InventoryHandler([]string{}, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	if !strings.Contains(buf.String(), "Sword") {
		t.Errorf("want output to contain %q, got %q", "Sword", buf.String())
	}
}

func TestDispatch_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	cmd := Command{Name: "fly", Args: []string{}}
	err = g.Registry.Dispatch(cmd, g)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	if !strings.Contains(buf.String(), "I don't know how to do that") {
		t.Errorf("want unknown command message, got %q", buf.String())
	}
}

func TestQuitHandler(t *testing.T) {
	var buf bytes.Buffer
	g, err := setupGame(&buf)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
	err = QuitHandler([]string{}, g)
	if err == nil {
		t.Fatal("expected error got nil")
	}

	if !errors.Is(err, ErrQuit) {
		t.Errorf("expected %v, got %v", ErrQuit, err)
	}
}

func TestRegisterDefaultHandlers(t *testing.T) {
	tests := []struct {
		name     string
		cmd      Command
		wantRoom string
	}{
		{
			name:     "n alias",
			cmd:      Command{Name: "n", Args: []string{}},
			wantRoom: "dining",
		},
		{
			name:     "north alias",
			cmd:      Command{Name: "north", Args: []string{}},
			wantRoom: "dining",
		},
		{
			name:     "s alias",
			cmd:      Command{Name: "s", Args: []string{}},
			wantRoom: "kitchen",
		},
		{
			name:     "south alias",
			cmd:      Command{Name: "south", Args: []string{}},
			wantRoom: "kitchen",
		},
		{
			name:     "w alias",
			cmd:      Command{Name: "w", Args: []string{}},
			wantRoom: "library",
		},
		{
			name:     "west alias",
			cmd:      Command{Name: "west", Args: []string{}},
			wantRoom: "library",
		},
		{
			name:     "e alias",
			cmd:      Command{Name: "e", Args: []string{}},
			wantRoom: "sport",
		},
		{
			name:     "east alias",
			cmd:      Command{Name: "east", Args: []string{}},
			wantRoom: "sport",
		},
		{
			name:     "u alias",
			cmd:      Command{Name: "u", Args: []string{}},
			wantRoom: "entrance",
		},
		{
			name:     "up alias",
			cmd:      Command{Name: "up", Args: []string{}},
			wantRoom: "entrance",
		},
		{
			name:     "d alias",
			cmd:      Command{Name: "d", Args: []string{}},
			wantRoom: "entrance",
		},
		{
			name:     "down alias",
			cmd:      Command{Name: "down", Args: []string{}},
			wantRoom: "entrance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			g, err := setupAliasWorld(&buf)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			err = g.Registry.Dispatch(tt.cmd, g)
			if err != nil {
				t.Fatalf("expected no error got: %v", err)
			}

			if g.Player.CurrentRoom() != tt.wantRoom {
				t.Errorf("want: %q, got: %q", tt.wantRoom, g.Player.CurrentRoom())
			}
		})
	}

}
