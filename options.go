package lastcmt

import (
	"fmt"
)

type Options struct {
	Owner  string `short:"o" required:"" help:"Owner name."`
	Repo   string `short:"r" required:"" help:"Repo name."`
	Number int    `short:"n" required:"" help:"Pull Request number."`
	Key    string `short:"k" default:"lastcmt" help:"Identification key."`
	Token  string `required:"" env:"GITHUB_TOKEN" help:"Auth token."`
}

func (options *Options) HTMLCommentID() string {
	return fmt.Sprintf("<!-- lastcmt: %s -->", options.Key)
}
