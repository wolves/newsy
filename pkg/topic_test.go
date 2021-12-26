package newsy

import "testing"

func Test_Topics_Match(t *testing.T) {
	t.Parallel()

	topics := Topics{"tdd", "go", "neovim"}

	tests := []struct {
		name  string
		query Topics
		exp   bool
	}{
		{name: "partial match", query: Topics{"neovim", "emacs"}, exp: true},
		{name: "complete match", query: Topics{"tdd", "go"}, exp: true},
		{name: "non-match", query: Topics{"vscode"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topics.Match(tt.query...)

			if got != tt.exp {
				t.Fatalf("expected %v, got %v", tt.exp, got)
			}
		})
	}
}
