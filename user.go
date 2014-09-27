package gl

import (
	"net/url"
)

const (
	users_url = "/users"
	user_url  = "/users/:id"
)

type User struct {
	Id            int    `json:"id,omitempty"`
	Username      string `json:"username,omitempty"`
	Email         string `json:"email,omitempty"`
	Name          string `json:"name,omitempty"`
	State         string `json:"state,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	Bio           string `json:"bio,omitempty"`
	Skype         string `json:"skype,omitempty"`
	LinkedIn      string `json:"linkedin,omitempty"`
	Twitter       string `json:"twitter,omitempty"`
	ExternUid     string `json:"extern_uid,omitempty"`
	Provider      string `json:"provider,omitempty"`
	ThemeId       int    `json:"theme_id,omitempty"`
	ColorSchemeId int    `json:"color_scheme_id,color_scheme_id"`
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

func (g *Client) CreateUser(email, username, password, name string,
	skype, linkedin, twitter, website *string,
	limit *int, externUid, provider, bio *string, admin, canCreateGroup *bool) (*User, error) {
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
	addBool(vals, "admin", admin)
	addBool(vals, "can_create_group", canCreateGroup)
	e := g.post(users_url, vals, &us)
	if e != nil {
		return nil, e
	}
	return &us, nil
}
