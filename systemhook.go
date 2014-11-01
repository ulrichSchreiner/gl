package gl

import (
	"net/url"
	"time"
)

const (
	systemhooks_url = "/hooks"
	systemhook_url  = "/hooks/:id"
)

type SystemHook struct {
	Id        int       `json:"id,omitempty"`
	Url       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at, omitempty"`
}
type SystemHooks []SystemHook

type SystemHookResult struct {
	EventName  string `json:"event_name,omitempty"`
	Name       string `json:"name,omitempty"`
	Path       string `json:"path,omitempty"`
	OwnerName  string `json:"owner_name,omitempty"`
	OwnerEmail string `json:"owner_email,omitempty"`
	ProjectId  int    `json:"project_id,omitempty"`
}

func (g *Client) SystemHooks(pg *Page) (SystemHooks, *Pagination, error) {
	var p SystemHooks
	pager, e := g.get(systemhooks_url, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) AllSystemHooks() (SystemHooks, error) {
	var r SystemHooks
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.SystemHooks(pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) AddSystemHook(u string) (*SystemHook, error) {
	vals := make(url.Values)
	vals.Set("url", u)
	var h SystemHook
	e := g.post(systemhooks_url, vals, &h)
	return &h, e
}

func (g *Client) TestSystemHook(hid int) (*SystemHookResult, error) {
	u := expandUrl(systemhook_url, map[string]interface{}{":id": hid})
	var r SystemHookResult
	_, e := g.get(u, nil, nil, &r)
	return &r, e
}

func (g *Client) DeleteSystemHook(hid int) (*SystemHook, error) {
	u := expandUrl(systemhook_url, map[string]interface{}{":id": hid})
	var r SystemHook
	e := g.delete(u, nil, &r)
	return &r, e
}
