package gl

import (
	"fmt"
	"net/url"
)

const (
	users_url       = "/users"
	user_url        = "/users/:id"
	curuser_url     = "/user"
	curuserkeys_url = "/user/keys"
	key_url         = "/user/keys/:id"
	userkeys_url    = "/users/:uid/keys"
	userkey_url     = "/users/:uid/keys/:id"
	session_url     = "/session"
)

type User struct {
	Id               int    `json:"id,omitempty"`
	Username         string `json:"username,omitempty"`
	Email            string `json:"email,omitempty"`
	Name             string `json:"name,omitempty"`
	State            State  `json:"state,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	Bio              string `json:"bio,omitempty"`
	Skype            string `json:"skype,omitempty"`
	LinkedIn         string `json:"linkedin,omitempty"`
	Twitter          string `json:"twitter,omitempty"`
	ExternUid        string `json:"extern_uid,omitempty"`
	Provider         string `json:"provider,omitempty"`
	ThemeId          int    `json:"theme_id,omitempty"`
	ColorSchemeId    int    `json:"color_scheme_id,omitempty"`
	PrivateToken     string `json:"private_token,omitempty"`
	Blocked          bool   `json:"blocked,omitempty"`
	IsAdmin          bool   `json:"is_admin,omitempty"`
	CanCreateGroup   bool   `json:"can_create_group,omitempty"`
	CanCreateTeam    bool   `json:"can_create_team,omitempty"`
	CanCreateProject bool   `json:"can_create_project,omitempty"`
}

type SshKey struct {
	Id    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Key   string `json:"key,omitempty"`
}

func (g *Client) users(usersurl string, pg *Page, parm url.Values) ([]User, *Pagination, error) {
	var r []User
	pager, e := g.get(usersurl, parm, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}

func (g *Client) allUsers(f fetchFunc) ([]User, error) {
	var r []User
	err := fetchAll(f, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) Users(pg *Page) ([]User, *Pagination, error) {
	return g.users(users_url, pg, nil)
}
func (g *Client) AllUsers() ([]User, error) {
	return g.allUsers(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Users(pg)
	})
}
func (g *Client) SearchUsers(query string, pg *Page) ([]User, *Pagination, error) {
	v := make(url.Values)
	v.Set("search", query)
	return g.users(users_url, pg, v)
}
func (g *Client) SearchAllUsers(query string) ([]User, error) {
	return g.allUsers(func(pg *Page) (interface{}, *Pagination, error) {
		return g.SearchUsers(query, pg)
	})
}

func (g *Client) GetUser(uid int) (*User, error) {
	var us User
	u := expandUrl(user_url, map[string]interface{}{":id": uid})
	_, e := g.get(u, nil, nil, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}
func (g *Client) CurrentUser() (*User, error) {
	var us User
	_, e := g.get(curuser_url, nil, nil, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}
func (g *Client) CreateUser(email, username, password, name string,
	skype, linkedin, twitter, website *string,
	limit *int, externUid, provider, bio *string, admin, canCreateGroup bool) (*User, error) {
	var us User
	vals := make(url.Values)
	vals.Set("email", email)
	vals.Set("password", password)
	vals.Set("username", username)
	vals.Set("name", name)
	addString(vals, "skype", skype)
	addString(vals, "linkedin", linkedin)
	addString(vals, "twitter", twitter)
	addString(vals, "website_url", website)
	addInt(vals, "projects_limit", limit)
	addString(vals, "extern_uid", externUid)
	addString(vals, "provider", provider)
	addString(vals, "bio", bio)
	vals.Set("admin", fmt.Sprintf("%v", admin))
	vals.Set("can_create_group", fmt.Sprintf("%v", canCreateGroup))
	e := g.post(users_url, vals, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}

func (g *Client) EditUser(uid int, email, username, password, name string,
	skype, linkedin, twitter, website *string,
	limit *int, externUid, provider, bio *string, admin, canCreateGroup bool) (*User, error) {
	var us User
	vals := make(url.Values)
	vals.Set("email", email)
	vals.Set("password", password)
	vals.Set("username", username)
	vals.Set("name", name)
	addString(vals, "skype", skype)
	addString(vals, "linkedin", linkedin)
	addString(vals, "twitter", twitter)
	addString(vals, "website_url", website)
	addInt(vals, "projects_limit", limit)
	addString(vals, "extern_uid", externUid)
	addString(vals, "provider", provider)
	addString(vals, "bio", bio)
	vals.Set("admin", fmt.Sprintf("%v", admin))
	vals.Set("can_create_group", fmt.Sprintf("%v", canCreateGroup))
	u := expandUrl(user_url, map[string]interface{}{":id": uid})
	e := g.put(u, vals, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}

func (g *Client) DeleteUser(uid int) (*User, error) {
	var us User
	u := expandUrl(user_url, map[string]interface{}{":id": uid})
	e := g.delete(u, nil, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}

func (g *Client) sshkeys(keysurl string, pg *Page, parm url.Values) ([]SshKey, *Pagination, error) {
	var r []SshKey
	pager, e := g.get(keysurl, parm, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}

func (g *Client) allsshkeys(f fetchFunc) ([]SshKey, error) {
	var r []SshKey
	err := fetchAll(f, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (g *Client) CurrentUserKeys(pg *Page) ([]SshKey, *Pagination, error) {
	return g.sshkeys(curuserkeys_url, pg, nil)
}
func (g *Client) AllCurrentUserKeys() ([]SshKey, error) {
	return g.allsshkeys(func(pg *Page) (interface{}, *Pagination, error) {
		return g.CurrentUserKeys(pg)
	})
}
func (g *Client) UserKeys(uid int, pg *Page) ([]SshKey, *Pagination, error) {
	u := expandUrl(userkeys_url, map[string]interface{}{":id": uid})
	return g.sshkeys(u, pg, nil)
}
func (g *Client) AllUserKeys(uid int) ([]SshKey, error) {
	return g.allsshkeys(func(pg *Page) (interface{}, *Pagination, error) {
		return g.UserKeys(uid, pg)
	})
}
func (g *Client) GetSshKey(kid int) (*SshKey, error) {
	var us SshKey
	u := expandUrl(key_url, map[string]interface{}{":id": kid})
	_, e := g.get(u, nil, nil, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}
func (g *Client) CreateCurrentUserSshKey(title, key string) (*SshKey, error) {
	var k SshKey
	vals := make(url.Values)
	vals.Set("title", title)
	vals.Set("key", key)
	e := g.post(curuserkeys_url, vals, &k)
	if e != nil {
		return nil, e
	}
	return &k, nil
}
func (g *Client) DeleteCurrentUserKey(kid int) (*SshKey, error) {
	var us SshKey
	u := expandUrl(key_url, map[string]interface{}{":id": kid})
	e := g.delete(u, nil, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}
func (g *Client) CreateSshKey(uid int, title, key string) (*SshKey, error) {
	u := expandUrl(userkeys_url, map[string]interface{}{":uid": uid})
	var k SshKey
	vals := make(url.Values)
	vals.Set("title", title)
	vals.Set("key", key)
	e := g.post(u, vals, &k)
	if e != nil {
		return nil, e
	}
	return &k, nil
}
func (g *Client) DeleteUserKey(uid, kid int) (*SshKey, error) {
	var us SshKey
	u := expandUrl(userkey_url, map[string]interface{}{":uid": uid, ":id": kid})
	e := g.delete(u, nil, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}
func (g *Client) Session(login string, email *string, password string) (*User, error) {
	var u User
	vals := make(url.Values)
	vals.Set("login", login)
	addString(vals, "email", email)
	vals.Set("password", password)
	e := g.post(session_url, vals, &u)
	if e != nil {
		return nil, e
	}
	return &u, nil
}
