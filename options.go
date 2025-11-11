package lastcmt

import (
	"fmt"
	"html"
)

type Options struct {
	Owner        string `short:"o" required:"" help:"Owner name."`
	Repo         string `short:"r" required:"" help:"Repo name."`
	Number       int    `short:"n" required:"" help:"Pull Request number."`
	Key          string `short:"k" default:"lastcmt" help:"Identification key."`
	MinimizeOnly bool   `short:"m" negatable:"" help:"Minimize only."`
	Token        string `required:"" env:"GITHUB_TOKEN" help:"Auth token."`
}

func (options *Options) HTMLCommentID() string {
	return fmt.Sprintf("<!-- lastcmt: %s -->", html.EscapeString(options.Key))
}
