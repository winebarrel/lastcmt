package lastcmt

import (
	"context"
	"fmt"
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

type issueOrPullRequestComment struct {
	id          githubv4.ID
	body        string
	isMinimized bool
	createdAt   time.Time
	author      string
}

type issueOrPullRequest struct {
	id       githubv4.ID
	url      string
	comments []issueOrPullRequestComment
}

func (client *Client) CommentWithMinimize(ctx context.Context, body string) (string, error) {
	login, err := client.getViewer(ctx)

	if err != nil {
		return "", err
	}

	pr, err := client.getIssueOrPullRequest(ctx)

	if err != nil {
		return "", err
	}

	url := pr.url

	if !client.MinimizeOnly {
		bodyWithID := client.HTMLCommentID() + "\n" + body
		url, err = client.createComment(ctx, pr.id, bodyWithID)

		if err != nil {
			return "", err
		}
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

func (client *Client) getIssueOrPullRequest(ctx context.Context) (*issueOrPullRequest, error) {
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
			IssueOrPullRequest struct {
				Issue struct {
					Id       githubv4.ID
					URL      string
					Comments struct {
						Nodes    []comment
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"comments(first: 100, after: $issueCursor)"`
				} `graphql:"... on Issue"`
				PullRequest struct {
					Id       githubv4.ID
					URL      string
					Comments struct {
						Nodes    []comment
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"comments(first: 100, after: $pullRequestCursor)"`
				} `graphql:"... on PullRequest"`
			} `graphql:"issueOrPullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]any{
		"owner":             githubv4.String(client.Owner),
		"repo":              githubv4.String(client.Repo),
		"number":            githubv4.Int(client.Number),
		"issueCursor":       (*githubv4.String)(nil),
		"pullRequestCursor": (*githubv4.String)(nil),
	}

	issueOrPR := &issueOrPullRequest{}
	var allComments []comment

	for {
		err := client.raw.Query(ctx, &q, variables)

		if err != nil {
			return nil, err
		}

		qIssueOrPR := q.Repository.IssueOrPullRequest

		if qIssueOrPR.Issue.Id != nil {
			issue := qIssueOrPR.Issue
			issueOrPR.id = issue.Id
			issueOrPR.url = issue.URL
			allComments = append(allComments, issue.Comments.Nodes...)

			if !issue.Comments.PageInfo.HasNextPage {
				break
			}

			variables["issueCursor"] = githubv4.NewString(issue.Comments.PageInfo.EndCursor)
		} else if qIssueOrPR.PullRequest.Id != nil {
			pr := qIssueOrPR.PullRequest
			issueOrPR.id = pr.Id
			issueOrPR.url = pr.URL
			allComments = append(allComments, pr.Comments.Nodes...)

			if !pr.Comments.PageInfo.HasNextPage {
				break
			}

			variables["pullRequestCursor"] = githubv4.NewString(pr.Comments.PageInfo.EndCursor)
		} else {
			return nil, fmt.Errorf("Could not resolve to an issue or pull request with the number of %d.", client.Number) //nolint:staticcheck
		}
	}

	for _, c := range allComments {
		issueOrPR.comments = append(issueOrPR.comments, issueOrPullRequestComment{
			id:          c.ID,
			body:        c.Body,
			isMinimized: c.IsMinimized,
			createdAt:   c.CreatedAt.Time,
			author:      c.Author.Login,
		})
	}

	slices.SortFunc(issueOrPR.comments, func(i, j issueOrPullRequestComment) int { return i.createdAt.Compare(j.createdAt) })

	return issueOrPR, nil
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
