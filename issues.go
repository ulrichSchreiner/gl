package gl

import (
	"net/url"
	"strings"
	"time"
)

type IssueStateEvent string

const (
	CloseIssue  = IssueStateEvent("close")
	ReopenIssue = IssueStateEvent("reopen")
)

const (
	issues_url        = "/issues"
	projectissues_url = "/projects/:id/issues"
	projectissue_url  = "/projects/:id/issues/:issue_id"
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

func (g *Client) issues(endpoint string, pid string, state *IssueStateEvent, lbls []string, pg *Page) (Issues, *Pagination, error) {
	var p Issues
	vals := make(url.Values)
	if state != nil {
		vals.Set("state", string(*state))
	}
	if lbls != nil {
		labels := strings.Join(lbls, ",")
		vals.Set("labels", labels)
	}

	u := expandUrl(endpoint, map[string]interface{}{":id": pid})
	pager, e := g.get(u, vals, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) ProjectIssues(pid string, state *IssueStateEvent, lbls []string, pg *Page) (Issues, *Pagination, error) {
	return g.issues(projectissues_url, pid, state, lbls, pg)
}
func (g *Client) AllProjectIssues(pid string, state *IssueStateEvent, lbls []string) (Issues, error) {
	var is Issues
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.ProjectIssues(pid, state, lbls, pg)
	}, &is)
	if err != nil {
		return nil, err
	}
	return is, nil
}

func (g *Client) Issues(state *IssueStateEvent, lbls []string, pg *Page) (Issues, *Pagination, error) {
	return g.issues(issues_url, "", state, lbls, pg)
}
func (g *Client) AllIssues(pid string, state *IssueStateEvent, lbls []string) (Issues, error) {
	var is Issues
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Issues(state, lbls, pg)
	}, &is)
	if err != nil {
		return nil, err
	}
	return is, nil
}

func (g *Client) Issue(pid string, iid int) (*Issue, error) {
	u := expandUrl(projectissue_url, map[string]interface{}{":id": pid, ":issue_id": iid})
	var i Issue
	_, e := g.get(u, nil, nil, &i)
	return &i, e
}

func (g *Client) CreateIssue(pid string, title string, desc *string, assignee *int, milestone *int, labels []string) (*Issue, error) {
	u := expandUrl(projectissues_url, map[string]interface{}{":id": pid})
	vals := make(url.Values)
	vals.Set("title", title)
	addString(vals, "description", desc)
	addInt(vals, "assignee_id", assignee)
	addInt(vals, "milestone_id", milestone)
	if labels != nil {
		lbls := strings.Join(labels, ",")
		vals.Set("labels", lbls)
	}
	var i Issue
	e := g.post(u, vals, &i)
	return &i, e
}

func (g *Client) UpdateIssue(pid string, iid int, title, description *string, assignee *int, milestone *int, labels []string, state *IssueStateEvent) (*Issue, error) {
	u := expandUrl(projectissue_url, map[string]interface{}{":id": pid, ":issue_id": iid})
	vals := make(url.Values)
	addString(vals, "title", title)
	addString(vals, "description", description)
	addInt(vals, "assignee_id", assignee)
	addInt(vals, "milestone_id", milestone)
	if labels != nil {
		lbls := strings.Join(labels, ",")
		vals.Set("labels", lbls)
	}
	if state != nil {
		vals.Set("state_event", string(*state))
	}
	var i Issue
	e := g.put(u, vals, &i)
	return &i, e
}
