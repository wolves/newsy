package cli

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type App struct {
	IO

	cmds map[string]Commander
	once sync.Once
}

func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) == 0 || args[0] == "-h" {
		return app.Usage(app.Stdout())
	}

	if err := app.init(); err != nil {
		return err
	}

	cmd, ok := app.cmds[args[0]]
	if !ok {
		return fmt.Errorf("command %q not found", args[0])
	}

	if ioCmd, ok := cmd.(IOCommander); ok {
		ioCmd.SetIO(app.IO)
	}

	fmt.Println("\n##############")
	fmt.Println("args:", args)
	fmt.Println("pwd:", pwd)
	fmt.Printf("##############\n\n")
	return cmd.Main(ctx, pwd, args[1:])
}

func (app *App) Usage(w io.Writer) error {
	if err := app.init(); err != nil {
		return err
	}

	fmt.Fprintln(w, "Usage: newsy <command> [options] [<args>...]")
	fmt.Fprintln(w, "---------------")

	for k, v := range app.cmds {
		uc, ok := v.(UsageCommander)
		if !ok {
			fmt.Fprintf(w, "Undefined command: %s:\n", k)
			continue
		}

		err := uc.Usage(app.Stdout())
		if err != nil {
			return nil
		}

		fmt.Fprintln(w)
	}

	// - list categories in backup
	// - backup file location
	// - # Articles in backup
	// - # Articles / Topic
	// - # Articles / Source

	return nil
}

func (app *App) init() error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}
	app.once.Do(func() {
		if app.cmds == nil || len(app.cmds) == 0 {
			app.cmds = map[string]Commander{
				"clear":  &ClearCmd{},
				"read":   &ReadCmd{},
				"stream": &StreamCmd{},
			}
		}
	})

	return nil
}

// Root Command
// - list categories
// - backup file location
// - # Articles in backup
// - # Articles / Topic
// - # Articles / Source

// Read Command

// Clear Command
