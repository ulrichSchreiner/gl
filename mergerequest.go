package gl

import (
	"fmt"
	"net/url"
)

type MergeState string
type MergeOrderBy string

const (
	AllMerges    = MergeState("all")
	OpenedMerges = MergeState("opened")
	ClosedMerges = MergeState("closed")
	MergedMerges = MergeState("merged")

	OrderByCreated = MergeOrderBy("created_at")
	OrderByUpdated = MergeOrderBy("updated_at")
)

const (
	mergerequests_url = "/projects/:id/merge_requests"
	mergerequest_url  = "/projects/:id/merge_request/:merge_request_id"
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
