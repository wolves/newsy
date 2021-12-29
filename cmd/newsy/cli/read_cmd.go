package cli

import (
	"context"
	"flag"
	"io"
)

var (
	_ IOCommander    = &ReadCmd{}
	_ FlagCommander  = &ReadCmd{}
	_ UsageCommander = &ReadCmd{}
)

type ReadCmd struct {
	IO

	Backup string
	AsJSON bool
	Output string

	flags *flag.FlagSet
}

func (cmd *ReadCmd) SetIO(oi IO) {
	cmd.IO = oi
}

func (cmd *ReadCmd) Flags() *flag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}

	flags := flag.NewFlagSet("read", flag.ContinueOnError)
	flags.BoolVar(&cmd.AsJSON, "j", false, "Prints the news stories in JSON format.")
	flags.StringVar(&cmd.Backup, "f", "newsy_db.json", "Location of article archive/backup")
	flags.StringVar(&cmd.Output, "o", "", "Specifies location for article output (Replaces command line output)")

	// o := cmd.Stdout()
	flags.SetOutput(cmd.Stdout())

	cmd.flags = flags

	return flags
}

func (cmd *ReadCmd) Main(ctx context.Context, pwd string, args []string) error {
	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()

	if len(args) == 0 || args[0] == "-h" {
		return cmd.Usage(cmd.Stdout())
	}
	return nil
}

func (cmd *ReadCmd) Usage(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.Usage()

	return nil
}
