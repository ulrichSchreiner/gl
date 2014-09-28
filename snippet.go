package gl

import (
	"net/url"
	"time"
)

const (
	snippets_url    = "/projects/:id/snippets"
	snippet_url     = "/projects/:id/snippets/:snippet_id"
	snippet_content = "/projects/:id/snippets/:snippet_id/raw"
)

type SnippetAuthor struct {
	Id       int       `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Name     string    `json:"name,omitempty"`
	State    string    `json:"state,omitempty"`
	Created  time.Time `json:"created_at,omitempty"`
}

type Snippet struct {
	Id       int            `json:"id,omitempty"`
	Title    string         `json:"title,omitempty"`
	FileName string         `json:"file_name,omitempty"`
	Expires  time.Time      `json:"expires_at,omitempty"`
	Updated  time.Time      `json:"updated_at,omitempty"`
	Created  time.Time      `json:"created_at,omitempty"`
	Author   *SnippetAuthor `json:"author,omitempty"`
}

func (g *Client) Snippets(id string, pg *Page) ([]Snippet, *Pagination, error) {
	var r []Snippet
	parm := make(url.Values)
	u := expandUrl(snippets_url, map[string]interface{}{":id": id})
	pager, e := g.get(u, parm, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}

func (g *Client) AllSnippets(pid string) ([]Snippet, error) {
	var r []Snippet
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Snippets(pid, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) GetSnippet(pid string, snipid int) (*Snippet, error) {
	var s Snippet
	u := expandUrl(snippet_url, map[string]interface{}{":id": pid, ":snippet_id": snipid})
	_, e := g.get(u, nil, nil, &s)
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (g *Client) CreateSnippet(id string, title, filename, code string) (*Snippet, error) {
	u := expandUrl(snippets_url, map[string]interface{}{":id": id})
	var s Snippet
	vals := make(url.Values)
	vals.Set("title", title)
	vals.Set("file_name", filename)
	vals.Set("code", code)
	e := g.post(u, vals, &s)
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (g *Client) EditSnippet(id string, snipid int, title, filename, code *string) (*Snippet, error) {
	u := expandUrl(snippet_url, map[string]interface{}{":id": id, ":snippet_id": snipid})
	var s Snippet
	vals := make(url.Values)
	addString(vals, "title", title)
	addString(vals, "file_name", filename)
	addString(vals, "code", code)
	e := g.put(u, vals, &s)
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (g *Client) DeleteSnippet(id string, snipid int) (*Snippet, error) {
	u := expandUrl(snippet_url, map[string]interface{}{":id": id, ":snippet_id": snipid})
	var s Snippet
	e := g.delete(u, nil, &s)
	if e != nil {
		return nil, e
	}
	return &s, nil
}
func (g *Client) SnippetContent(id string, snipid int) ([]byte, error) {
	u := expandUrl(snippet_content, map[string]interface{}{":id": id, ":snippet_id": snipid})
	buf, _, err := g.httpexecute("GET", u, nil, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
