package main

import (
	"context"
	"fmt"
	"io"
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
		BodyFile kong.FileContentFlag `arg:"" optional:"" type:"filecontent" xor:"stdin" help:"Comment body file. If not specified, read from stdin."`
		Version  kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	args := os.Args[1:]

	if _, err := parser.Parse(args); err != nil {
		parser.FatalIfErrorf(err)
	}

	if len(args) == 1 && !cli.MinimizeOnly {
		if stdin, err := io.ReadAll(os.Stdin); err != nil {
			parser.FatalIfErrorf(err)
		} else {
			cli.BodyFile = stdin
		}
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
