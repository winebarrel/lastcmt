package lastcmt

import (
	"fmt"
	"html"
)

type Options struct {
	Owner        string `short:"o" required:"" env:"GITHUB_OWNER" help:"Owner name."`
	Repo         string `short:"r" required:"" env:"GITHUB_REPO" help:"Repo name."`
	Number       int    `short:"n" required:"" help:"Issue/Pull Request number."`
	Key          string `short:"k" default:"lastcmt" help:"Identification key."`
	MinimizeOnly bool   `short:"m" negatable:"" help:"Minimize only."`
	Token        string `required:"" env:"GITHUB_TOKEN" help:"Auth token."`
}

func (options *Options) HTMLCommentID() string {
	return fmt.Sprintf("<!-- lastcmt: %s -->", html.EscapeString(options.Key))
}
