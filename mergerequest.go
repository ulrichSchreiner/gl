package gl

import (
	"fmt"
	"net/url"
)

type MergeState string
type MergeEvent string
type MergeOrderBy string

const (
	AllMerges    = MergeState("all")
	OpenedMerges = MergeState("opened")
	ClosedMerges = MergeState("closed")
	MergedMerges = MergeState("merged")

	CloseMerge  MergeEvent = "close"
	ReopenMerge            = "reopen"
	MergeMerge             = "merge"

	OrderByCreated = MergeOrderBy("created_at")
	OrderByUpdated = MergeOrderBy("updated_at")
)

const (
	mergerequests_url = "/projects/:id/merge_requests"
	mergerequest_url  = "/projects/:id/merge_request/:merge_request_id"
	merge_url         = "/projects/:id/merge_request/:merge_request_id/merge"
	commentmerge_url  = "/projects/:id/merge_request/:merge_request_id/comment"
	commentsmerge_url = "/projects/:id/merge_request/:merge_request_id/comments"
)

type MergeRequest struct {
	Id           int        `json:"id,omitempty"`
	Iid          int        `json:"iid,omitempty"`
	TargetBranch string     `json:"target_branch,omitempty"`
	SourceBranch string     `json:"source_branch,omitempty"`
	ProjectId    int        `json:"project_id,omitempty"`
	Title        string     `json:"title,omitempty"`
	State        MergeState `json:"state,omitempty"`
	Upvotes      int        `json:"upvotes,omitempty"`
	Downvotes    int        `json:"downvotes,omitempty"`
	Author       *User      `json:"author,omitempty"`
	Assignee     *User      `json:"assignee,omitempty"`
	Description  string     `json:"description,omitempty"`
}

type MergeComment struct {
	Author *User  `json:"author,omitempty"`
	Note   string `json:"note,omitempty"`
}

func (g *Client) MergeRequests(id string, state *MergeState, orderBy *MergeOrderBy, asc *bool, pg *Page) ([]MergeRequest, *Pagination, error) {
	var r []MergeRequest
	parm := make(url.Values)
	if state != nil {
		parm.Set("state", fmt.Sprintf("%s", *state))
	}
	if orderBy != nil {
		parm.Set("order_by", fmt.Sprintf("%s", *orderBy))
	}
	if asc != nil {
		if *asc {
			parm.Set("sort", "asc")
		} else {
			parm.Set("sort", "desc")
		}
	}
	u := expandUrl(mergerequests_url, map[string]interface{}{":id": id})
	pager, e := g.get(u, parm, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}

func (g *Client) AllMergeRequests(pid string, state *MergeState, orderBy *MergeOrderBy, asc *bool) ([]MergeRequest, error) {
	var r []MergeRequest
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.MergeRequests(pid, state, orderBy, asc, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) GetMergeRequest(pid string, mrid int) (*MergeRequest, error) {
	var s MergeRequest
	u := expandUrl(mergerequest_url, map[string]interface{}{":id": pid, ":merge_request_id": mrid})
	_, e := g.get(u, nil, nil, &s)
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (g *Client) CreateMergeRequest(pid string, sbranch, tbranch string, assignee *int, title string, targetproject *int) (*MergeRequest, error) {
	vals := make(url.Values)
	vals.Set("source_branch", sbranch)
	vals.Set("target_branch", tbranch)
	addInt(vals, "assignee_id", assignee)
	vals.Set("title", title)
	addInt(vals, "target_project_id", targetproject)

	u := expandUrl(mergerequest_url, map[string]interface{}{":id": pid})

	var m MergeRequest
	err := g.post(u, vals, &m)
	return &m, err
}

func (g *Client) UpdateMergeRequest(pid string, mid int, sbranch, tbranch string, assignee int, title string, state MergeState) (*MergeRequest, error) {
	vals := make(url.Values)
	vals.Set("source_branch", sbranch)
	vals.Set("target_branch", tbranch)
	addInt(vals, "assignee_id", &assignee)
	vals.Set("title", title)
	vals.Set("state_event", string(state))

	u := expandUrl(mergerequests_url, map[string]interface{}{":id": pid, ":merge_request_id": mid})

	var m MergeRequest
	err := g.put(u, vals, &m)
	return &m, err
}

func (g *Client) AcceptMerge(pid string, mid int, msg *string) (*MergeRequest, error) {
	vals := make(url.Values)
	addString(vals, "merge_commit_message", msg)
	u := expandUrl(merge_url, map[string]interface{}{":id": pid, ":merge_request_id": mid})

	var m MergeRequest
	err := g.put(u, vals, &m)
	return &m, err
}

func (g *Client) CommentMerge(pid string, mid int, msg *string) (*MergeComment, error) {
	vals := make(url.Values)
	addString(vals, "note", msg)
	u := expandUrl(commentmerge_url, map[string]interface{}{":id": pid, ":merge_request_id": mid})

	var m MergeComment
	err := g.put(u, vals, &m)
	return &m, err
}

func (g *Client) MergeComments(id string, mid int, pg *Page) ([]MergeComment, *Pagination, error) {
	var r []MergeComment
	u := expandUrl(commentsmerge_url, map[string]interface{}{":id": id, ":merge_request_id": mid})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}

func (g *Client) AllMergeComments(pid string, mid int) ([]MergeComment, error) {
	var r []MergeComment
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.MergeComments(pid, mid, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
