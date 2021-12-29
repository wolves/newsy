package cli

import (
	"context"
	"flag"
	"io"
)

type Commander interface {
	Main(ctx context.Context, pwd string, args []string) error
}

type IOCommander interface {
	Commander
	SetIO(oi IO)
}

type UsageCommander interface {
	Commander
	Usage(w io.Writer) error
}

type FlagCommander interface {
	Commander
	Flags() *flag.FlagSet
}
