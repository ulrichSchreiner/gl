package gl

import (
	"net/url"
	"time"
)

const (
	tags_url     = "/projects/:id/repository/tags"
	tree_url     = "/projects/:id/repository/tree"
	file_content = "/projects/:id/repository/blobs/:sha"
	blob_content = "/projects/:id/repository/raw_blobs/:sha"
	archive_url  = "/projects/:id/repository/archive"
)

type CommitParent struct {
	Id string `json:"id,omitempty"`
}
type RepoCommit struct {
	Id        string         `json:"id,omitempty"`
	Message   string         `json:"message,omitempty"`
	Tree      string         `json:"tree,omitempty"`
	Author    *PersonData    `json:"author,omitempty"`
	Committer *PersonData    `json:"committer,omitempty"`
	Authored  time.Time      `json:"authored_date,omitempty"`
	Committed time.Time      `json:"committed_date,omitempty"`
	Parents   []CommitParent `json:"parents,omitempty"`
}
type inlineRepoCommit struct {
	Id             string         `json:"id,omitempty"`
	Message        string         `json:"message,omitempty"`
	Tree           string         `json:"tree,omitempty"`
	AuthorName     string         `json:"author_name,omitempty"`
	AuthorEMail    string         `json:"author_email,omitempty"`
	CommitterName  string         `json:"committer_name,omitempty"`
	CommitterEMail string         `json:"committer_email,omitempty"`
	Authored       time.Time      `json:"authored_date,omitempty"`
	Committed      time.Time      `json:"committed_date,omitempty"`
	Parents        []CommitParent `json:"parents,omitempty"`
}
type NamedCommit struct {
	Name      string      `json:"name,omitempty"`
	Commit    *RepoCommit `json:"commit,omitempty"`
	Protected bool        `json:"protected,omitempty"`
}
type namedCommit struct {
	Name      string            `json:"name,omitempty"`
	Commit    *inlineRepoCommit `json:"commit,omitempty"`
	Protected bool              `json:"protected,omitempty"`
}

type Branch struct {
	NamedCommit
}
type Tag struct {
	NamedCommit
}

type RepositoryEntry struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Mode string `json:"mode,omitempty"`
}

func (g *Client) Branches(id *int, nsname *string, pg *Page) ([]Branch, *Pagination, error) {
	if err := checkName(id, nsname); err != nil {
		return nil, nil, err
	}
	var r []Branch
	u := expandUrl(branches_url, map[string]interface{}{":id": idname(id, nsname)})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllBranches(id *int, nsname *string) ([]Branch, error) {
	if err := checkName(id, nsname); err != nil {
		return nil, err
	}
	var b []Branch
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Branches(id, nsname, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (g *Client) Branch(id *int, nsname *string, branch string) (*Branch, error) {
	if err := checkName(id, nsname); err != nil {
		return nil, err
	}
	var b Branch
	u := expandUrl(branch_url, map[string]interface{}{":id": idname(id, nsname), ":branch": branch})
	_, e := g.get(u, nil, nil, &b)
	if e != nil {
		return nil, e
	}
	return &b, nil
}

func (g *Client) protectBranch(id *int, nsname *string, branch string, command string) (*Branch, error) {
	if err := checkName(id, nsname); err != nil {
		return nil, err
	}
	var b Branch
	u := expandUrl(branch_url, map[string]interface{}{":id": idname(id, nsname), ":branch": branch})
	u = u + command
	if e := g.put(u, nil, &b); e != nil {
		return nil, e
	}
	return &b, nil
}

func (g *Client) ProtectBranch(id *int, nsname *string, branch string) (*Branch, error) {
	return g.protectBranch(id, nsname, branch, "/protect")
}
func (g *Client) UnprotectBranch(id *int, nsname *string, branch string) (*Branch, error) {
	return g.protectBranch(id, nsname, branch, "/unprotect")
}

func (g *Client) Tags(pid int, pg *Page) ([]Tag, *Pagination, error) {
	var r []Tag
	u := expandUrl(tags_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllTags(id int) ([]Tag, error) {
	var b []Tag
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Tags(id, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (g *Client) CreateTag(id int, name, ref string, msg *string) (*Tag, error) {
	u := expandUrl(tags_url, map[string]interface{}{":id": id})
	var s namedCommit
	vals := make(url.Values)
	vals.Set("tag_name", name)
	vals.Set("ref", ref)
	addString(vals, "message", msg)
	e := g.post(u, vals, &s)
	if e != nil {
		return nil, e
	}
	var t Tag
	t.Name = s.Name
	t.Protected = s.Protected
	var c RepoCommit
	t.Commit = &c
	c.Author = &PersonData{Name: s.Commit.AuthorName, EMail: s.Commit.AuthorEMail}
	c.Committer = &PersonData{Name: s.Commit.CommitterName, EMail: s.Commit.CommitterEMail}
	c.Id = s.Commit.Id
	c.Message = s.Commit.Message
	c.Tree = s.Commit.Tree
	c.Authored = s.Commit.Authored
	c.Committed = s.Commit.Committed
	c.Parents = s.Commit.Parents
	return &t, nil
}

func (g *Client) RepoEntries(id int, path, ref *string, pg *Page) ([]RepositoryEntry, *Pagination, error) {
	var r []RepositoryEntry
	u := expandUrl(tree_url, map[string]interface{}{":id": id})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllRepoEntries(id int, path, ref *string) ([]RepositoryEntry, error) {
	var b []RepositoryEntry
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.RepoEntries(id, path, ref, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (g *Client) RawFileContent(id int, sha, filepath string) ([]byte, error) {
	u := expandUrl(file_content, map[string]interface{}{":id": id, ":sha": sha})
	p := make(url.Values)
	p.Set("filepath", filepath)
	buf, _, err := g.httpexecute("GET", u, p, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (g *Client) RawBlobContent(id int, sha string) ([]byte, error) {
	u := expandUrl(blob_content, map[string]interface{}{":id": id, ":sha": sha})
	buf, _, err := g.httpexecute("GET", u, nil, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (g *Client) Archive(id int, sha *string) ([]byte, error) {
	u := expandUrl(archive_url, map[string]interface{}{":id": id})
	v := make(url.Values)
	addString(v, "sha", sha)
	buf, _, err := g.httpexecute("GET", u, v, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
