package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

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

	f, err := os.Open(fmt.Sprintf("%s/%s", pwd, cmd.Backup))
	if err != nil {
		return err
	}
	defer f.Close()

	if !cmd.AsJSON {
		sess, err := newsy.Restore(f)
		if err != nil {
			return err
		}
		arts, err := fetchArticles(sess.Articles, args...)
		var resp string
		for _, a := range arts {
			resp += fmt.Sprintf("\n%s\n", a.String())
		}

		if cmd.Output != "" {
			fp := pwd + "/" + cmd.Output
			f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.Write([]byte(resp)); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(cmd.Stdout(), "%s", resp)
		}
		return err
	}

	if _, err = io.Copy(cmd.Stdout(), f); err != nil {
		return err
	}
	return err
}

func (cmd *ReadCmd) Usage(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.Usage()

	return nil
}

func fetchArticles(as newsy.Articles, args ...string) (newsy.Articles, error) {
	matches := newsy.Articles{}

	for _, sid := range args {
		id, err := strconv.Atoi(sid)
		if err != nil {
			return nil, err
		}
		for _, a := range as {
			if a.ID == newsy.ArticleID(id) {
				matches = append(matches, a)
			}
		}
	}

	return matches, nil
}
