package goforth

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCommand Command
	}{
		{
			name:        "empty string",
			input:       "",
			wantCommand: Command{},
		},
		{
			name:        "whitespaces string",
			input:       "           ",
			wantCommand: Command{},
		},
		{
			name:        "one word",
			input:       "go",
			wantCommand: Command{Name: "go", Args: []string{}},
		},
		{
			name:        "multiple words",
			input:       "pick up sword",
			wantCommand: Command{Name: "pick", Args: []string{"up", "sword"}},
		},
		{
			name:        "upper case words",
			input:       "GO NORTH",
			wantCommand: Command{Name: "go", Args: []string{"north"}},
		},
		{
			name:        "mixed case words",
			input:       "gO noRtH",
			wantCommand: Command{Name: "go", Args: []string{"north"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Parse(tt.input)
			if result.Name != tt.wantCommand.Name {
				t.Errorf("want: %q, got: %q", tt.wantCommand.Name, result.Name)
			}
		})
	}

}
