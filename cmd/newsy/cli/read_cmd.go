package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	newsy "github.com/wolves/newsy/pkg"
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
	flags.StringVar(&cmd.Output, "o", "", "Specifies location for command output (Replaces command line output)")

	o := cmd.Stdout()
	flags.SetOutput(o)

	cmd.flags = flags

	return flags
}

func (cmd *ReadCmd) Main(ctx context.Context, pwd string, args []string) error {
	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()

	if len(args) == 0 {
		return fmt.Errorf("please provide 1 or more article ids")
	}

	// fmt.Printf("Args: %v", args)

	f, err := os.Open(fmt.Sprintf("%s/%s", pwd, cmd.Backup))
	if err != nil {
		return err
	}
	defer f.Close()

	sess, err := newsy.Restore(f)
	if err != nil {
		return err
	}
	// fmt.Printf("Session: %+v", sess)

	for _, a := range sess.Articles {
		fmt.Fprintf(cmd.Stdout(), "\n%s\n", a.String())
		// _, err = io.Copy(cmd.Out, f)
	}

	return err
}

func (cmd *ReadCmd) Usage(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.Usage()

	return nil
}
