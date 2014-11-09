package gl

import (
	"net/url"
	"time"
)

const (
	tags_url         = "/projects/:id/repository/tags"
	tree_url         = "/projects/:id/repository/tree"
	file_content     = "/projects/:id/repository/blobs/:sha"
	blob_content     = "/projects/:id/repository/raw_blobs/:sha"
	archive_url      = "/projects/:id/repository/archive"
	compare_url      = "/projects/:id/repository/compare"
	contributors_url = "/projects/:id/repository/contributors"
	readfile_url     = "/projects/:id/repository/files"
	commits_url      = "/projects/:id/repository/commits"
	commit_url       = "/projects/:id/repository/commits/:sha"
	commitdiff_url   = "/projects/:id/repository/commits/:sha/diff"
)

type CommitParent struct {
	Id string `json:"id,omitempty"`
}
type CommitEx struct {
	Id        string         `json:"id,omitempty"`
	Message   string         `json:"message,omitempty"`
	Tree      string         `json:"tree,omitempty"`
	Author    *PersonData    `json:"author,omitempty"`
	Committer *PersonData    `json:"committer,omitempty"`
	Authored  time.Time      `json:"authored_date,omitempty"`
	Committed time.Time      `json:"committed_date,omitempty"`
	Parents   []CommitParent `json:"parents,omitempty"`
}

type SimpleCommit struct {
	ShortId     string    `json:"short_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	AuthorName  string    `json:"author_name,omitempty"`
	AuthorEMail string    `json:"author_email,omitempty"`
	Created     time.Time `json:"created_at,omitempty"`
}

type Commit struct {
	Id             string         `json:"id,omitempty"`
	ShortId        string         `json:"short_id,omitempty"`
	Message        string         `json:"message,omitempty"`
	Tree           string         `json:"tree,omitempty"`
	AuthorName     string         `json:"author_name,omitempty"`
	AuthorEMail    string         `json:"author_email,omitempty"`
	CommitterName  string         `json:"committer_name,omitempty"`
	CommitterEMail string         `json:"committer_email,omitempty"`
	Created        time.Time      `json:"created_date,omitempty"`
	Authored       time.Time      `json:"authored_date,omitempty"`
	Committed      time.Time      `json:"committed_date,omitempty"`
	Parents        []CommitParent `json:"parents,omitempty"`
}
type NamedCommitEx struct {
	Name      string    `json:"name,omitempty"`
	Commit    *CommitEx `json:"commit,omitempty"`
	Protected bool      `json:"protected,omitempty"`
}
type namedCommit struct {
	Name      string  `json:"name,omitempty"`
	Commit    *Commit `json:"commit,omitempty"`
	Protected bool    `json:"protected,omitempty"`
}

type Branch struct {
	NamedCommitEx
}
type TagListEntry struct {
	NamedCommitEx
}
type Tag struct {
	Commit
}

type RepositoryEntry struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Mode string `json:"mode,omitempty"`
}

type Diff struct {
	NewPath     string  `json:"new_path,omitempty"`
	OldPath     string  `json:"old_path,omitempty"`
	Amode       *string `json:"a_mode,omitempty"`
	Bmode       *string `json:"b_mode,omitempty"`
	Diff        string  `json:"diff,omitempty"`
	NewFile     bool    `json:"new_file,omitempty"`
	RenamedFile bool    `json:"renamed_file,omitempty"`
	DeletedFile bool    `json:"deleted_file,omitempty"`
}

type Comparison struct {
	Commit   *SimpleCommit  `json:"commit,omitempty"`
	Commits  []SimpleCommit `json:"commits,omitempty"`
	Diffs    []Diff         `json:"diffs,omitempty"`
	Timemout bool           `json:"compare_timeout,omitempty"`
	SameRef  bool           `json:"compare_same_ref,omitempty"`
}

type Contributor struct {
	Name      string `json:"name,omitempty"`
	EMail     string `json:"email,omitempty"`
	Commits   int    `json:"commits,omitempty"`
	Additions int    `json:"additions,omitempty"`
	Deletions int    `json:"deletions,omitempty"`
}

type RepoFile struct {
	Branch   string `json:"branch_name,omitempty"`
	Name     string `json:"file_name,omitempty"`
	Path     string `json:"file_path,omitempty"`
	Size     int    `json:"size,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	Content  string `json:"content,omitempty"`
	Ref      string `json:"ref,omitempty"`
	BlobId   string `json:"blob_id,omitempty"`
	CommitId string `json:"commit_id,omitempty"`
}

func (g *Client) Branches(id string, pg *Page) ([]Branch, *Pagination, error) {
	var r []Branch
	u := expandUrl(branches_url, map[string]interface{}{":id": id})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllBranches(id string) ([]Branch, error) {
	var b []Branch
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Branches(id, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (g *Client) Branch(id string, branch string) (*Branch, error) {
	var b Branch
	u := expandUrl(branch_url, map[string]interface{}{":id": id, ":branch": branch})
	_, e := g.get(u, nil, nil, &b)
	if e != nil {
		return nil, e
	}
	return &b, nil
}

func (g *Client) protectBranch(id string, branch string, command string) (*Branch, error) {
	var b Branch
	u := expandUrl(branch_url, map[string]interface{}{":id": id, ":branch": branch})
	u = u + command
	if e := g.put(u, nil, &b); e != nil {
		return nil, e
	}
	return &b, nil
}

func (g *Client) ProtectBranch(id string, branch string) (*Branch, error) {
	return g.protectBranch(id, branch, "/protect")
}
func (g *Client) UnprotectBranch(id string, branch string) (*Branch, error) {
	return g.protectBranch(id, branch, "/unprotect")
}

func (g *Client) Tags(pid string, pg *Page) ([]TagListEntry, *Pagination, error) {
	var r []TagListEntry
	u := expandUrl(tags_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllTags(id string) ([]TagListEntry, error) {
	var b []TagListEntry
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Tags(id, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (g *Client) CreateTag(id string, name, ref string, msg *string) (*Tag, error) {
	u := expandUrl(tags_url, map[string]interface{}{":id": id})
	var s Tag
	vals := make(url.Values)
	vals.Set("tag_name", name)
	vals.Set("ref", ref)
	addString(vals, "message", msg)
	e := g.post(u, vals, &s)
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (g *Client) RepoEntries(id string, path, ref *string, pg *Page) ([]RepositoryEntry, *Pagination, error) {
	var r []RepositoryEntry
	u := expandUrl(tree_url, map[string]interface{}{":id": id})
	vals := make(url.Values)
	addString(vals, "path", path)
	addString(vals, "ref_name", ref)
	pager, e := g.get(u, vals, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllRepoEntries(id string, path, ref *string) ([]RepositoryEntry, error) {
	var b []RepositoryEntry
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.RepoEntries(id, path, ref, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (g *Client) RawFileContent(id string, sha, filepath string) ([]byte, error) {
	u := expandUrl(file_content, map[string]interface{}{":id": id, ":sha": sha})
	p := make(url.Values)
	p.Set("filepath", filepath)
	buf, _, err := g.httpexecute("GET", u, p, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (g *Client) RawBlobContent(id string, sha string) ([]byte, error) {
	u := expandUrl(blob_content, map[string]interface{}{":id": id, ":sha": sha})
	buf, _, err := g.httpexecute("GET", u, nil, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (g *Client) Archive(id string, sha *string) ([]byte, error) {
	u := expandUrl(archive_url, map[string]interface{}{":id": id})
	v := make(url.Values)
	addString(v, "sha", sha)
	buf, _, err := g.httpexecute("GET", u, v, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (g *Client) Compare(id string, from, to string) (*Comparison, error) {
	u := expandUrl(compare_url, map[string]interface{}{":id": id})
	v := make(url.Values)
	v.Set("from", from)
	v.Set("to", to)
	var c Comparison
	_, e := g.get(u, v, nil, &c)
	if e != nil {
		return nil, e
	}
	return &c, nil
}
func (g *Client) Contributors(id string, pg *Page) ([]Contributor, *Pagination, error) {
	var r []Contributor
	u := expandUrl(contributors_url, map[string]interface{}{":id": id})
	pager, e := g.get(u, nil, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllContributors(id string) ([]Contributor, error) {
	var b []Contributor
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Contributors(id, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (g *Client) ReadFile(id, filepath, ref string) (*RepoFile, error) {
	var b RepoFile
	u := expandUrl(readfile_url, map[string]interface{}{":id": id})
	v := make(url.Values)
	v.Set("file_path", filepath)
	v.Set("ref", ref)
	_, e := g.get(u, v, nil, &b)
	if e != nil {
		return nil, e
	}
	return &b, nil
}
func (g *Client) CreateFile(id, filepath, branch, commitmsg, content string, encoding string) (*RepoFile, error) {
	return g.changeFile(true, id, filepath, branch, commitmsg, content, encoding)
}
func (g *Client) UpdateFile(id, filepath, branch, commitmsg, content string, encoding string) (*RepoFile, error) {
	return g.changeFile(false, id, filepath, branch, commitmsg, content, encoding)
}
func (g *Client) changeFile(ispost bool, id, filepath, branch, commitmsg, content string, encoding string) (*RepoFile, error) {
	var b RepoFile
	u := expandUrl(readfile_url, map[string]interface{}{":id": id})
	v := make(url.Values)
	v.Set("file_path", filepath)
	v.Set("branch_name", branch)
	v.Set("encoding", encoding)
	v.Set("content", content)
	v.Set("commit_message", commitmsg)
	var e error
	if ispost {
		e = g.post(u, v, &b)
	} else {
		e = g.put(u, v, &b)
	}
	if e != nil {
		return nil, e
	}
	return &b, nil
}
func (g *Client) DeleteFile(id, filepath, branch, commitmsg string) (*RepoFile, error) {
	var b RepoFile
	u := expandUrl(readfile_url, map[string]interface{}{":id": id})
	v := make(url.Values)
	v.Set("file_path", filepath)
	v.Set("branch_name", branch)
	v.Set("commit_message", commitmsg)
	e := g.delete(u, v, &b)
	if e != nil {
		return nil, e
	}
	return &b, nil
}

func (g *Client) Commits(id string, ref *string, pg *Page) ([]Commit, *Pagination, error) {
	var r []Commit
	u := expandUrl(commits_url, map[string]interface{}{":id": id})
	vals := make(url.Values)
	addString(vals, "ref_name", ref)
	pager, e := g.get(u, vals, pg, &r)
	if e != nil {
		return nil, nil, e
	}
	return r, pager, nil
}
func (g *Client) AllCommits(id string, ref *string) ([]Commit, error) {
	var b []Commit
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Commits(id, ref, pg)
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (g *Client) ReadCommit(id, sha string) (*Commit, error) {
	var b Commit
	u := expandUrl(commit_url, map[string]interface{}{":id": id, ":sha": sha})
	_, e := g.get(u, nil, nil, &b)
	if e != nil {
		return nil, e
	}
	return &b, nil
}
func (g *Client) ReadDiff(id, sha string) (*Diff, error) {
	var b Diff
	u := expandUrl(commitdiff_url, map[string]interface{}{":id": id, ":sha": sha})
	_, e := g.get(u, nil, nil, &b)
	if e != nil {
		return nil, e
	}
	return &b, nil
}
