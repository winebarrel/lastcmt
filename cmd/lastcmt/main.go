package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/lastcmt"
)

var version string

func init() {
	log.SetFlags(0)
}

func parseArgs() (string, *lastcmt.Options) {
	var cli struct {
		lastcmt.Options
		BodyFile kong.FileContentFlag `arg:"" optional:"" type:"filecontent" help:"Comment body file. '-' is accepted for stdin."`
		Version  kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	args := os.Args[1:]
	_, err := parser.Parse(args)

	parser.FatalIfErrorf(err)

	if !cli.MinimizeOnly && len(args) == 0 {
		parser.FatalIfErrorf(errors.New(`expected "<body-file>"`))
	}
	return string(cli.BodyFile), &cli.Options
}

func main() {
	body, options := parseArgs()
	ctx := context.Background()
	client := lastcmt.NewClient(ctx, options)
	url, err := client.CommentWithMinimize(ctx, body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)
}
