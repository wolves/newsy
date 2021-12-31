package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func Test_Cli(t *testing.T) {
	t.Run("root cmd", func(t *testing.T) {
		app := &App{
			cmds: map[string]Commander{},
		}

		bb := &bytes.Buffer{}
		app.Out = bb
		args := []string{}

		err := app.Main(context.Background(), "test_pwd", args)
		assertNoError(t, err)

		act := bb.String()[0:5]

		exp := "Usage"

		if !strings.Contains(act, exp) {
			t.Fatalf("expected: %v, got: %v", exp, act)
		}
	})

	t.Run("sub-command routing", func(t *testing.T) {
		app := &App{
			cmds: map[string]Commander{},
		}

		args := []string{"stream"}
		err := app.Main(context.Background(), "test_pwd", args)
		assertNoError(t, err)
	})
}

// Helpers
func assertNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertStringEquality(t testing.TB, got, exp string) {
	t.Helper()

	if got != exp {
		t.Fatalf("expected: %v, got: %v", exp, got)
	}
}
