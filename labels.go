package gl

import (
	"net/url"
)

const (
	labels_url = "/projects/:id/labels"
)

type Label struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}
type Labels []Label

func (g *Client) Labels(pid string, pg *Page) (Labels, *Pagination, error) {
	var p Labels
	u := expandUrl(labels_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) AllLabels(pid string) (Labels, error) {
	var r Labels
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Labels(pid, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) CreateLabel(pid, name, color string) (*Label, error) {
	u := expandUrl(labels_url, map[string]interface{}{":id": pid})
	vals := make(url.Values)
	vals.Set("name", name)
	vals.Set("color", color)
	var l Label
	e := g.post(u, vals, &l)
	return &l, e
}

func (g *Client) DeleteLabel(pid, name string) (*Label, error) {
	u := expandUrl(labels_url, map[string]interface{}{":id": pid})
	vals := make(url.Values)
	vals.Set("name", name)
	var l Label
	e := g.delete(u, vals, &l)
	return &l, e
}

func (g *Client) UpdateLabel(pid, name string, newname, color *string) (*Label, error) {
	u := expandUrl(labels_url, map[string]interface{}{":id": pid})
	vals := make(url.Values)
	vals.Set("name", name)
	addString(vals, "color", color)
	addString(vals, "new_name", newname)
	var l Label
	e := g.put(u, vals, &l)
	return &l, e
}
