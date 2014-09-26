package gl

import (
	"fmt"
	"github.com/spacemonkeygo/errors"
	"net/url"
	"strconv"
	"time"
)

const (
	projects_url       = "/projects"
	projects_all_url   = "/projects/all"
	projects_owned_url = "/projects/owned"
	project_url        = "/projects/:id"
	project_events_url = "/projects/:id/events"
	userprojects_url   = "/projects/user/:user_id"
	members_url        = "/projects/:id/members"
	member_url         = "/projects/:id/members/:user_id"
	hooks_url          = "/projects/:id/hooks"
)

type MemberState string

var InvalidParam = errors.NewClass("invalid parameter")

const (
	MemberActive MemberState = "active"
)

type Permission struct {
	Access       AccessLevel       `json:"access_level,omitempty"`
	Notification NotificationLevel `json:"notification_level,omitempty"`
}

type Permissions struct {
	Project *Permission `json:"project_access,omitempty"`
	Group   *Permission `json:"group_access,omitempty"`
}

// The member type in gitlab.
type Member struct {
	Id       int         `json:"id,omitempty"`
	Username string      `json:"username,omitempty"`
	EMail    string      `json:"email,omitempty"`
	Name     string      `json:"name,omitempty"`
	State    MemberState `json:"state,omitempty"`
	Created  time.Time   `json:"created_at,omitempty"`
	Access   AccessLevel `json:"access_level,omitempty"`
}
type Members []Member

// The namespace type in gitlab.
type Namespace struct {
	Id          int       `json:"id,omitempty"`
	Created     time.Time `json:"created_at,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	OwnerId     int       `json:"owner_id,omitempty"`
	Path        string    `json:"path,omitempty"`
	Updated     time.Time `json:"updated_at,omitempty"`
}

// the project type in gitlab.
type Project struct {
	Id                   int             `json:"id,omitempty"`
	Description          string          `json:"description,omitempty"`
	DefaultBranch        string          `json:"default_branch,omitempty"`
	Public               bool            `json:"public,omitempty"`
	Visibility           VisibilityLevel `json:"visibility,omitempty"`
	SshRepoUrl           string          `json:"ssh_url_to_repo,omitempty"`
	HttpRepoUrl          string          `json:"http_url_to_repo,omitempty"`
	WebUrl               string          `json:"web_url,omitempty"`
	Owner                *Member         `json:"owner,omitempty"`
	Name                 string          `json:"name,omitempty"`
	NameWithSpaces       string          `json:"name_with_spaces,omitempty"`
	Path                 string          `json:"path,omitempty"`
	PathWithSpaces       string          `json:"path_with_spaces,omitempty"`
	IssuesEnabled        bool            `json:"issues_enabled,omitempty"`
	MergeRequestsEnabled bool            `json:"merge_requests_enabled,omitempty"`
	WikiEnabled          bool            `json:"wiki_enabled,omitempty"`
	SnippetsEnabled      bool            `json:"snippets_enabled,omitempty"`
	Created              time.Time       `json:"created_at, omitempty"`
	LastActivity         time.Time       `json:"last_activity_at, omitempty"`
	Archived             bool            `json:"archived, omitempty"`
	Permissions          Permissions     `json:"permissions,omitempty"`
}
type Projects []Project

type EventData struct {
	Before       string      `json:"before,omitempty"`
	After        string      `json:"after,omitempty"`
	Ref          string      `json:"ref,omitempty"`
	UserId       int         `json:"user_id,omitempty"`
	UserName     string      `json:"user_name,omitempty"`
	Repository   *Repository `json:"repository,omitempty"`
	Commits      []Commit    `json:"commits,omitempty"`
	TotalCommits int         `json:"total_commits_count,omitempty"`
}

type Repository struct {
	Name        string `json:"name,omitempty"`
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Homepage    string `json:"homepage,omitempty"`
}

type Commit struct {
	Id         string    `json:"id,omitempty"`
	Message    string    `json:"message,omitempty"`
	Timestampe time.Time `json:"timestamp,omitempty"`
	URL        string    `json:"url,omitempty"`
	Author     *Author   `json:"author,omitempty"`
}
type Author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
type Event struct {
	ProjectId   int        `json:"project_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	ActionName  string     `json:"action_name,omitempty"`
	TargetId    int        `json:"target_id,omitempty"`
	TargetType  string     `json:"target_type,omitempty"`
	AuthorId    int        `json:"author_id,omitempty"`
	Data        *EventData `json:"data,omitempty"`
	TargetTitle string     `json:"target_title,omitempty"`
}
type Events []Event

type Hook struct {
	Id                  int       `json:"id,omitempty"`
	Url                 string    `json:"url,omitempty"`
	ProjectId           int       `json:"project_id,omitempty"`
	PushEvents          bool      `json:"push_events,omitempty"`
	IssuesEvents        bool      `json:"issues_events,omitempty"`
	MergeRequestsEvents bool      `json:"merge_requests_events,omitempty"`
	CreatedAt           time.Time `json:"created_at, omitempty"`
}

func (g *Client) projects(purl string, pg *Page) (Projects, *Pagination, error) {
	var p Projects
	pager, e := g.get(purl, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) allProjects(f fetchFunc) (Projects, error) {
	var p Projects
	err := fetchAll(f, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (g *Client) VisibleProjects(pg *Page) (Projects, *Pagination, error) {
	return g.projects(projects_url, pg)
}
func (g *Client) Projects(pg *Page) (Projects, *Pagination, error) {
	return g.projects(projects_all_url, pg)
}
func (g *Client) OwnedProjects(pg *Page) (Projects, *Pagination, error) {
	return g.projects(projects_owned_url, pg)
}

func (g *Client) AllVisibleProjects() (Projects, error) {
	return g.allProjects(func(pg *Page) (interface{}, *Pagination, error) {
		return g.VisibleProjects(pg)
	})
}
func (g *Client) AllOwnedProjects() (Projects, error) {
	return g.allProjects(func(pg *Page) (interface{}, *Pagination, error) {
		return g.OwnedProjects(pg)
	})
}
func (g *Client) AllProjects() (Projects, error) {
	return g.allProjects(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Projects(pg)
	})
}

func (g *Client) Project(id int) (*Project, error) {
	var p Project
	u := expandUrl(project_url, map[string]interface{}{":id": id})
	_, e := g.get(u, nil, nil, &p)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func (g *Client) Events(pid int, pg *Page) (Events, *Pagination, error) {
	var p Events
	u := expandUrl(project_events_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}
func (g *Client) AllEvents(pid int) (Events, error) {
	var ev Events
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Events(pid, pg)
	}, &ev)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

func (g *Client) CreateProject(name string, path *string, nsid *int, description *string,
	issuesEnabled, mergeRQenabled, wikiEnabled, snippetsEnabled *bool,
	public *bool, vis *VisibilityLevel, importUrl *string) (*Project, error) {
	return g.createProject(projects_url, nil, name, path, nsid, description, nil, issuesEnabled, mergeRQenabled, wikiEnabled, snippetsEnabled,
		public, vis, importUrl)
}
func (g *Client) CreateUserProject(name string, uid int, description, defaultbranch *string,
	issuesEnabled, mergeRQenabled, wikiEnabled, snippetsEnabled *bool,
	public *bool, vis *VisibilityLevel, importUrl *string) (*Project, error) {
	return g.createProject(userprojects_url, map[string]interface{}{":user_id": uid}, name, nil, nil, description, defaultbranch,
		issuesEnabled, mergeRQenabled, wikiEnabled, snippetsEnabled,
		public, vis, importUrl)
}
func (g *Client) createProject(purl string, urlparms map[string]interface{}, name string, path *string, nsid *int, description, defbranch *string,
	issuesEnabled, mergeRQenabled, wikiEnabled, snippetsEnabled *bool,
	public *bool, vis *VisibilityLevel, importUrl *string) (*Project, error) {
	vals := make(url.Values)
	vals.Set("name", name)
	addString(vals, "path", path)
	addInt(vals, "namespace_id", nsid)
	addString(vals, "default_branch", defbranch)
	addString(vals, "description", description)
	addBool(vals, "issues_enabled", issuesEnabled)
	addBool(vals, "merge_requests_enabled", mergeRQenabled)
	addBool(vals, "wiki_enabled", wikiEnabled)
	addBool(vals, "snippets_enabled", snippetsEnabled)
	addBool(vals, "public", public)
	if vis != nil {
		v := int(*vis)
		addInt(vals, "visibility_level", &v)
	}
	addString(vals, "import_url", importUrl)

	if urlparms != nil {
		purl = expandUrl(purl, urlparms)
	}
	var p Project
	err := g.post(purl, vals, &p)
	return &p, err
}

func (g *Client) RemoveProject(id int) (*Project, error) {
	var p Project
	u := expandUrl(project_url, map[string]interface{}{":id": id})
	e := g.delete(u, nil, &p)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func (g *Client) AllTeamMembers(id *int, nsname *string, query *string) (Members, error) {
	var p Members
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.TeamMembers(id, nsname, query, pg)
	}, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func idname(id *int, nam *string) string {
	sid := ""
	if id != nil {
		sid = strconv.Itoa(*id)
	} else {
		sid = *nam
	}
	return sid
}

func (g *Client) TeamMembers(id *int, nsname *string, query *string, pg *Page) ([]Member, *Pagination, error) {
	if id == nil && nsname == nil {
		return nil, nil, InvalidParam.New("teamid or name must be given")
	}

	u := expandUrl(members_url, map[string]interface{}{":id": idname(id, nsname)})
	var p Members

	pager, e := g.get(u, nil, nil, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) TeamMember(pid, uid int) (*Member, error) {
	var p Member
	u := expandUrl(member_url, map[string]interface{}{":id": pid, ":user_id": uid})
	_, e := g.get(u, nil, nil, &p)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func (g *Client) AddTeamMember(id *int, nsname *string, uid int, level AccessLevel) (*Member, error) {
	if id == nil && nsname == nil {
		return nil, InvalidParam.New("projectid or name must be given")
	}
	var p Member
	u := expandUrl(members_url, map[string]interface{}{":id": idname(id, nsname)})
	vals := make(url.Values)
	vals.Set("access_level", fmt.Sprintf("%d", level))
	vals.Set("user_id", fmt.Sprintf("%d", uid))
	e := g.post(u, vals, &p)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func (g *Client) EditTeamMember(id *int, nsname *string, uid int, level AccessLevel) (*Member, error) {
	if id == nil && nsname == nil {
		return nil, InvalidParam.New("projectid or name must be given")
	}
	u := expandUrl(member_url, map[string]interface{}{":id": idname(id, nsname), ":user_id": uid})
	var p Member
	vals := make(url.Values)
	vals.Set("access_level", fmt.Sprintf("%d", level))
	e := g.put(u, vals, &p)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func (g *Client) DeleteTeamMember(id *int, nsname *string, uid int) (*Member, error) {
	if id == nil && nsname == nil {
		return nil, InvalidParam.New("projectid or name must be given")
	}
	u := expandUrl(member_url, map[string]interface{}{":id": idname(id, nsname), ":user_id": uid})
	var p Member
	e := g.delete(u, nil, &p)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

func (g *Client) Hooks(id *int, nsname *string, pg *Page) ([]Hook, *Pagination, error) {
	var p []Hook
	if id == nil && nsname == nil {
		return nil, nil, InvalidParam.New("projectid or name must be given")
	}
	u := expandUrl(hooks_url, map[string]interface{}{":id": idname(id, nsname)})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}
func (g *Client) AllHooks(id *int, nsname *string) ([]Hook, error) {
	var h []Hook
	if id == nil && nsname == nil {
		return nil, InvalidParam.New("projectid or name must be given")
	}
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Hooks(id, nsname, pg)
	}, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}
