package lastcmt

import (
	"context"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	rBotSuffix = regexp.MustCompile(`\[bot\]$`)
)

type Client struct {
	*Options
	raw *githubv4.Client
}

func NewClient(ctx context.Context, options *Options) *Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: options.Token})
	hc := oauth2.NewClient(ctx, src)
	rawCli := githubv4.NewClient(hc)
	client := &Client{Options: options, raw: rawCli}
	return client
}

type pullRequestComment struct {
	id          githubv4.ID
	body        string
	isMinimized bool
	createdAt   time.Time
	author      string
}

type pullRequest struct {
	id       githubv4.ID
	comments []pullRequestComment
}

func (client *Client) CommentWithMinimize(ctx context.Context, body string) (string, error) {
	login, err := client.getViewer(ctx)

	if err != nil {
		return "", err
	}

	pr, err := client.getPullRequest(ctx)

	if err != nil {
		return "", err
	}

	bodyWithID := client.HTMLCommentID() + "\n" + body
	url, err := client.createComment(ctx, pr.id, bodyWithID)

	if err != nil {
		return "", err
	}

	for _, c := range pr.comments {
		if c.author == login && !c.isMinimized && strings.Contains(c.body, client.HTMLCommentID()) {
			err := client.minimizeComment(ctx, c.id)

			if err != nil {
				return "", err
			}
		}
	}

	return url, nil
}

func (client *Client) getViewer(ctx context.Context) (string, error) {
	var q struct {
		Viewer struct {
			Login string
		}
	}

	err := client.raw.Query(ctx, &q, nil)

	if err != nil {
		return "", err
	}

	login := q.Viewer.Login
	login = rBotSuffix.ReplaceAllString(login, "")

	return login, nil
}

func (client *Client) getPullRequest(ctx context.Context) (*pullRequest, error) {
	type comment struct {
		ID          githubv4.ID
		Body        string
		IsMinimized bool
		CreatedAt   githubv4.DateTime
		Author      struct {
			Login string
		}
	}

	var q struct {
		Repository struct {
			PullRequest struct {
				Id       githubv4.ID
				Comments struct {
					Nodes    []comment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $cursor)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]any{
		"owner":  githubv4.String(client.Owner),
		"repo":   githubv4.String(client.Repo),
		"number": githubv4.Int(client.Number),
		"cursor": (*githubv4.String)(nil),
	}

	pr := &pullRequest{}
	var allComments []comment

	for {
		err := client.raw.Query(ctx, &q, variables)

		if err != nil {
			return nil, err
		}

		pr.id = q.Repository.PullRequest.Id
		allComments = append(allComments, q.Repository.PullRequest.Comments.Nodes...)

		if !q.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}

		variables["commentsCursor"] = githubv4.NewString(q.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}

	for _, c := range allComments {
		pr.comments = append(pr.comments, pullRequestComment{
			id:          c.ID,
			body:        c.Body,
			isMinimized: c.IsMinimized,
			createdAt:   c.CreatedAt.Time,
			author:      c.Author.Login,
		})
	}

	slices.SortFunc(pr.comments, func(i, j pullRequestComment) int { return i.createdAt.Compare(j.createdAt) })

	return pr, nil
}

func (client *Client) createComment(ctx context.Context, prID githubv4.ID, body string) (string, error) {
	var m struct {
		AddComment struct {
			CommentEdge struct {
				Node struct {
					URL string
				}
			}
		} `graphql:"addComment(input: $input)"`
	}

	input := githubv4.AddCommentInput{
		SubjectID: prID,
		Body:      githubv4.String(body),
	}

	err := client.raw.Mutate(ctx, &m, input, nil)

	if err != nil {
		return "", err
	}

	return m.AddComment.CommentEdge.Node.URL, nil
}

func (client *Client) minimizeComment(ctx context.Context, commentID githubv4.ID) error {
	var m struct {
		MinimizeComment struct {
			MinimizedComment struct {
				IsMinimized bool
			}
		} `graphql:"minimizeComment(input: $input)"`
	}

	input := githubv4.MinimizeCommentInput{
		SubjectID:  commentID,
		Classifier: githubv4.ReportedContentClassifiersOutdated,
	}

	err := client.raw.Mutate(ctx, &m, input, nil)

	if err != nil {
		return err
	}

	return nil
}
