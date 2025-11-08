package main

import (
	"context"
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
		Body    kong.FileContentFlag `arg:"" required:"" type:"filecontent" help:"Comment body file. '-' is accepted for stdin."`
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return string(cli.Body), &cli.Options
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
