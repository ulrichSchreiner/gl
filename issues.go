package gl

import (
	"time"
)

type IssueStateEvent string

const (
	CloseIssue  IssueStateEvent = "close"
	ReopenIssue IssueStateEvent = "reopen"
)

const (
	issues_url = "/projects/:id/issues"
	issue_url  = "/projects/:id/issues/:issue_id"
)

type Issue struct {
	Id          int        `json:"id,omitempty"`
	Iid         int        `json:"iid,omitempty"`
	ProjectId   int        `json:"project_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	State       string     `json:"state,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	Milestone   *Milestone `json:"milestone,omitempty"`
	Assignee    *User      `json:"assignee,omitempty"`
	Author      *User      `json:"author,omitempty"`
}

type Issues []Issue

func (g *Client) Issues(pid int, pg *Page) (Issues, *Pagination, error) {
	var p Issues
	u := expandUrl(issues_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}
func (g *Client) AllIssues(pid int) (Issues, error) {
	var is Issues
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Issues(pid, pg)
	}, &is)
	if err != nil {
		return nil, err
	}
	return is, nil
}
