package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	// pubsub "github.com/wolves/newsy/pkg"
)

type App struct {
	IO
	Commands map[string]Commander
}

type Commander interface {
	Main(ctx context.Context, pwd string, args []string) error
}

func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) == 0 {
		return app.Usage(app.Stdout())
	}

	if app.Commands == nil {
		app.Commands = map[string]Commander{
			// "search": &SearchCmd{
			// 	Name: "search",
			// },
		}
	}

	cmd, ok := app.Commands[args[0]]
	if !ok {
		return fmt.Errorf("command %q not found", args[0])
	}

	return cmd.Main(ctx, pwd, args[1:])

	// m := pubsub.NewManager()
	// fmt.Println("app.Main")
	// fmt.Println("args:", args)
	// fmt.Println("pwd:", pwd)
	// fmt.Printf("mgr: %+v", m)
}

func (app *App) Usage(w io.Writer) error {
	fmt.Fprintln(w, "Usage: newsy <command> [options] [<args>...]")
	fmt.Fprintln(w, "---------------")
	// - list categories
	// - backup file location
	// - # Articles in backup
	// - # Articles / Topic
	// - # Articles / Source

	return nil
}

// Root Command
// - list categories
// - backup file location
// - # Articles in backup
// - # Articles / Topic
// - # Articles / Source

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func (oi IO) Stdout() io.Writer {
	if oi.Out == nil {
		return os.Stdout
	}

	return oi.Out
}

func (oi IO) Stderr() io.Writer {
	if oi.Err == nil {
		return os.Stderr
	}

	return oi.Err
}

func (oi IO) Stdin() io.Reader {
	if oi.In == nil {
		return os.Stdin
	}

	return oi.In
}

// Stream Command

// Read Command

// Clear Command
