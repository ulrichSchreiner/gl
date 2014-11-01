package gl

import (
	"net/url"
	"time"
)

const (
	deploykeys_url = "/projects/:id/keys"
	deploykey_url  = "/projects/:id/keys/:key_id"
)

type DeployKey struct {
	Id      int        `json:"id,omitempty"`
	Title   string     `json:"title,omitempty"`
	Key     string     `json:"key,omitempty"`
	Created *time.Time `json:"created,omitempty"`
}
type DeployKeys []DeployKey

func (g *Client) DeployKeys(pid string, pg *Page) (DeployKeys, *Pagination, error) {
	var p DeployKeys
	u := expandUrl(deploykeys_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) AllDeployKeys(pid string) (DeployKeys, error) {
	var r DeployKeys
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.DeployKeys(pid, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) DeployKey(pid string, kid int) (*DeployKey, error) {
	u := expandUrl(deploykey_url, map[string]interface{}{":id": pid, ":key_id": kid})
	var i DeployKey
	_, e := g.get(u, nil, nil, &i)
	return &i, e
}

func (g *Client) AddKey(pid string, title, key string) (*DeployKey, error) {
	u := expandUrl(deploykeys_url, map[string]interface{}{":id": pid})
	v := make(url.Values)
	v.Set("title", title)
	v.Set("key", key)
	var k DeployKey
	e := g.post(u, v, &k)
	return &k, e
}

func (g *Client) RemoveKey(pid string, kid int) (*DeployKey, error) {
	u := expandUrl(deploykey_url, map[string]interface{}{":id": pid, ":key_id": kid})
	var i DeployKey
	e := g.delete(u, nil, &i)
	return &i, e

}
