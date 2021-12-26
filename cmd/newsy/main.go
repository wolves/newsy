package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/wolves/newsy/cmd/newsy/cli"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	// app := &cli.App{}
	app := &cli.App{}

	err = app.Main(ctx, pwd, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// fmt.Println("pwd:", pwd)

	<-ctx.Done()
}
