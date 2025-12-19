package lastcmt

import (
	"fmt"
	"html"
	"strings"
)

type Repo string

func (repo Repo) split() []string {
	return strings.Split(string(repo), "/")
}

func (repo Repo) Validate() error {
	parts := repo.split()

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return nil
	}

	return fmt.Errorf(`expected the "OWNER/REPO" format, got "%s"`, repo)
}

func (repo Repo) Owner() string {
	parts := repo.split()
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func (repo Repo) Name() string {
	parts := repo.split()
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

type Options struct {
	Repo         Repo   `short:"R" required:"" env:"GITHUB_REPOSITORY" help:"OWNER/REPO"`
	Number       int    `arg:"" required:"" help:"Issue/Pull Request number."`
	Key          string `short:"k" default:"lastcmt" help:"Identification key."`
	MinimizeOnly bool   `short:"m" negatable:"" xor:"stdin" help:"Minimize only."`
	Token        string `required:"" env:"GH_TOKEN,GITHUB_TOKEN" help:"Auth token."`
}

func (options *Options) HTMLCommentID() string {
	return fmt.Sprintf("<!-- lastcmt: %s -->", html.EscapeString(options.Key))
}
