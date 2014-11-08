package gl

import (
	"net/url"
	"time"
)

const (
	groups_url       = "/groups"
	group_url        = "/groups/:id"
	projectgroup_url = "/groups/:id/projects/:project_id"
	groupmembers_url = "/groups/:id/members"
	groupmember_url  = "/groups/:id/members/:user_id"
)

type Group struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Path    string `json:"path,omitempty"`
	OwnerId int    `json:"owner_id, omitempty"`
}
type Groups []Group

type GroupMember struct {
	Id       int         `json:"id,omitempty"`
	Username string      `json:"username,omitempty"`
	Email    string      `json:"email,omitempty"`
	Name     string      `json:"name,omitempty"`
	State    State       `json:"state,omitempty"`
	Created  *time.Time  `json:"created_at,omitempty"`
	Access   AccessLevel `json:"access_level,omitempty"`
}
type GroupMembers []GroupMember

func (g *Client) Groups(pg *Page) (Groups, *Pagination, error) {
	var p Groups
	pager, e := g.get(groups_url, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) AllGroups() (Groups, error) {
	var r Groups
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Groups(pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) Group(gid int) (*Group, error) {
	u := expandUrl(group_url, map[string]interface{}{":id": gid})
	var i Group
	_, e := g.get(u, nil, nil, &i)
	return &i, e
}

func (g *Client) AddGroup(name, path string) (*Group, error) {
	v := make(url.Values)
	v.Set("name", name)
	v.Set("path", path)

	var gr Group
	e := g.post(groups_url, v, &gr)
	return &gr, e
}

func (g *Client) TransferProjectToGroup(gid, pid int) (*Group, error) {
	v := make(url.Values)
	u := expandUrl(projectgroup_url, map[string]interface{}{":id": gid, ":project_id": pid})
	var gr Group
	e := g.post(u, v, &gr)
	return &gr, e
}

func (g *Client) DeleteGroup(gid int) (*Group, error) {
	u := expandUrl(group_url, map[string]interface{}{":id": gid})
	var r Group
	e := g.delete(u, nil, &r)
	return &r, e
}

func (g *Client) GroupMembers(gid int, pg *Page) (GroupMembers, *Pagination, error) {
	var p GroupMembers
	u := expandUrl(groupmembers_url, map[string]interface{}{":id": gid})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) AllGroupMembers(gid int) (GroupMembers, error) {
	var r GroupMembers
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.GroupMembers(gid, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) AddGroupMember(gid, uid int, level AccessLevel) (*GroupMember, error) {
	u := expandUrl(groupmembers_url, map[string]interface{}{":id": gid})
	v := make(url.Values)
	addInt(v, "user_id", &uid)
	l := int(level)
	addInt(v, "access_level", &l)

	var gr GroupMember
	e := g.post(u, v, &gr)
	return &gr, e
}

func (g *Client) DeleteGroupMember(gid, uid int) error {
	u := expandUrl(groupmember_url, map[string]interface{}{":id": gid, ":user_id": uid})
	return g.delete(u, nil, nil)
}
