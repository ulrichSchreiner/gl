package gl

import (
	"net/url"
	"time"
)

const (
	issuenotes_url   = "/projects/:id/issues/:issue_id/notes"
	issuenote_url    = "/projects/:id/issues/:issue_id/notes/:note_id"
	snippetnotes_url = "/projects/:id/snippets/:snippet_id/notes"
	snippetnote_url  = "/projects/:id/snippets/:snippet_id/notes/:note_id"
	mergenotes_url   = "/projects/:id/merge_requests/:merge_request_id/notes"
	mergenote_url    = "/projects/:id/merge_requests/:merge_request_id/notes/:note_id"
)

type Note struct {
	Id         int        `json:"id,omitempty"`
	Body       string     `json:"body,omitempty"`
	Attachment []byte     `json:"attachment,omitempty"`
	Title      string     `json:"title,omitempty"`
	Filename   string     `json:"file_name,omitempty"`
	Author     *User      `json:"author,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}
type Notes []Note

func (g *Client) notes(endpoint string, pid string, ntype string, nkey int, pg *Page) (Notes, *Pagination, error) {
	var p Notes

	u := expandUrl(endpoint, map[string]interface{}{":id": pid, ntype: nkey})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) note(endpoint string, pid string, noteid int, ntype string, nkey int) (*Note, error) {
	u := expandUrl(endpoint, map[string]interface{}{":id": pid, ":note_id": noteid, ntype: nkey})
	var n Note
	_, e := g.get(u, nil, nil, &n)
	return &n, e
}

func (g *Client) IssueNotes(pid string, iid int, pg *Page) (Notes, *Pagination, error) {
	return g.notes(issuenotes_url, pid, ":issue_id", iid, pg)
}
func (g *Client) AllIssueNotes(pid string, iid int) (Notes, error) {
	var n Notes
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.IssueNotes(pid, iid, pg)
	}, &n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (g *Client) IssueNote(pid string, iid int, nid int) (*Note, error) {
	return g.note(issuenote_url, pid, nid, ":issue_id", iid)
}
func (g *Client) CreateIssueNote(pid string, iid int, body string) (*Note, error) {
	u := expandUrl(issuenotes_url, map[string]interface{}{":id": pid, ":issue_id": iid})
	v := make(url.Values)
	v.Set("body", body)
	var n Note
	e := g.post(u, v, &n)
	return &n, e
}

func (g *Client) SnippetNotes(pid string, sid int, pg *Page) (Notes, *Pagination, error) {
	return g.notes(snippetnotes_url, pid, ":snippet_id", sid, pg)
}
func (g *Client) AllSnippetNotes(pid string, sid int) (Notes, error) {
	var n Notes
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.SnippetNotes(pid, sid, pg)
	}, &n)
	if err != nil {
		return nil, err
	}
	return n, nil
}
func (g *Client) SnippetNote(pid string, sid int, nid int) (*Note, error) {
	return g.note(snippetnote_url, pid, nid, ":snippet_id", sid)
}
func (g *Client) CreateSnippetNote(pid string, sid int, body string) (*Note, error) {
	u := expandUrl(snippetnotes_url, map[string]interface{}{":id": pid, ":snippet_id": sid})
	v := make(url.Values)
	v.Set("body", body)
	var n Note
	e := g.post(u, v, &n)
	return &n, e
}

func (g *Client) MergeNotes(pid string, mid int, pg *Page) (Notes, *Pagination, error) {
	return g.notes(mergenotes_url, pid, ":merge_request_id", mid, pg)
}
func (g *Client) AllMergeNotes(pid string, mid int) (Notes, error) {
	var n Notes
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.MergeNotes(pid, mid, pg)
	}, &n)
	if err != nil {
		return nil, err
	}
	return n, nil
}
func (g *Client) MergeNote(pid string, mid int, nid int) (*Note, error) {
	return g.note(mergenote_url, pid, nid, ":merge_request_id", mid)
}
func (g *Client) CreateMergeNote(pid string, mid int, body string) (*Note, error) {
	u := expandUrl(mergenotes_url, map[string]interface{}{":id": pid, ":merge_request_id": mid})
	v := make(url.Values)
	v.Set("body", body)
	var n Note
	e := g.post(u, v, &n)
	return &n, e
}
