package cli

import (
	"context"
	"flag"
	"io"
)

var (
	_ IOCommander    = &ClearCmd{}
	_ UsageCommander = &ClearCmd{}
	_ FlagCommander  = &ClearCmd{}
)

type ClearCmd struct {
	IO

	Backup string

	flags *flag.FlagSet
}

func (cmd *ClearCmd) SetIO(oi IO) {
	cmd.IO = oi
}

func (cmd *ClearCmd) Flags() *flag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}

	flags := flag.NewFlagSet("clear", flag.ContinueOnError)
	flags.StringVar(&cmd.Backup, "f", "newsy_db.json", "Location of Newsy article archive/backup")

	cmd.flags = flags

	return flags
}

func (cmd *ClearCmd) Main(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return cmd.Usage(cmd.Stdout())
	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	return nil
}

func (cmd *ClearCmd) Usage(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.Usage()

	return nil
}
