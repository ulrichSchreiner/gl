package gl

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
	"time"
)

func TestProjects(t *testing.T) {
	Convey("Creating a project with only a name", t, func() {
		name := "testproject"

		h := th(true, func(v url.Values) (interface{}, error, int) {
			var res Project
			res.Name = name
			return &res, nil, 200
		})
		srv, cl := StubHandler(h)
		defer srv.Close()
		p, _ := cl.CreateProject(name, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		Convey("it must be a post to the correct url", func() {
			So(h.method, ShouldEqual, "POST")
			So(h.path, ShouldEqual, projects_url)
		})
		Convey("the mapped URL parameters must be correct", func() {
			So(h.values.Get("name"), ShouldEqual, name)
			So(h.values, hasnot,
				"path",
				"namespace_id",
				"description",
				"merge_requests_enabled",
				"issues_enabled", "wiki_enabled", "snippets_enabled",
				"public",
				"visibility_level", "import_url")
		})
		Convey("and the result should be correct unmarshalled", func() {
			So(p.Name, ShouldEqual, name)
		})
	})
	Convey("Creating a full filled project", t, func() {
		name := "testproject"
		namespaceid := 1
		path := "p"
		desc := "desc"
		pub, merge, iss, wiki, snipp := true, true, true, true, true
		vis := Internal
		iurl := "import-url"
		h := th(true, func(v url.Values) (interface{}, error, int) {
			var res Project
			res.Name = name
			res.Path = path
			res.Description = desc
			res.Public = pub
			res.MergeRequestsEnabled = merge
			res.IssuesEnabled = merge
			res.WikiEnabled = wiki
			res.SnippetsEnabled = snipp
			res.Visibility = vis
			return &res, nil, 200
		})

		srv, cl := StubHandler(h)
		defer srv.Close()
		p, _ := cl.CreateProject(name, &path, &namespaceid, &desc, &iss, &merge, &wiki, &snipp, &pub, &vis, &iurl)
		Convey("the mapped URL parameters must be correct", func() {
			So(h.get("name"), ShouldEqual, name)
			So(h.values, has,
				"path",
				"namespace_id",
				"description",
				"merge_requests_enabled",
				"issues_enabled", "wiki_enabled", "snippets_enabled",
				"public",
				"visibility_level", "import_url")
		})
		Convey("and the result should be correct unmarshalled", func() {
			So(p.Name, ShouldEqual, name)
			So(p.Path, ShouldEqual, path)
			So(p.Description, ShouldEqual, desc)
			So(p.Public, ShouldEqual, pub)
			So(p.MergeRequestsEnabled, ShouldEqual, merge)
			So(p.IssuesEnabled, ShouldEqual, iss)
			So(p.WikiEnabled, ShouldEqual, wiki)
			So(p.SnippetsEnabled, ShouldEqual, snipp)
			So(p.Visibility, ShouldEqual, vis)
		})
	})

	Convey("Creating a userproject with a name and a defaultbranch", t, func() {
		name := "testproject"
		user := 5
		branch := "testbranch"
		h := th(true, func(v url.Values) (interface{}, error, int) {
			var res Project
			res.Name = name
			res.DefaultBranch = branch
			return &res, nil, 200
		})
		srv, cl := StubHandler(h)
		defer srv.Close()
		p, _ := cl.CreateUserProject(name, user, nil, &branch,
			nil, nil, nil, nil,
			nil, nil, nil)

		Convey("it must be a post to the correct url", func() {
			So(h.method, ShouldEqual, "POST")
			So(h.path, ShouldEqual, fmt.Sprintf("/projects/user/%d", user))
		})
		Convey("the mapped URL parameters must be correct", func() {
			So(h.get("name"), ShouldEqual, name)
			So(h.get("default_branch"), ShouldEqual, branch)
			So(h.values, hasnot,
				"description",
				"merge_requests_enabled",
				"issues_enabled", "wiki_enabled", "snippets_enabled",
				"public",
				"visibility_level", "import_url")
		})
		Convey("and the result should be correct unmarshalled", func() {
			So(p.Name, ShouldEqual, name)
			So(p.DefaultBranch, ShouldEqual, branch)
		})
	})

	Convey("removing a given project", t, func() {
		pid := 54
		h := th(true, func(v url.Values) (interface{}, error, int) {
			var p Project
			p.Name = "to be deleted"
			return &p, nil, 200
		})
		srv, cl := StubHandler(h)
		defer srv.Close()
		p, _ := cl.RemoveProject(pid)
		Convey("should invoke a DELETE on the correct url", func() {
			So(h.method, ShouldEqual, "DELETE")
			So(h.path, ShouldEqual, "/projects/54")
			Convey("and the result should be the deleted project", func() {
				So(p.Name, ShouldEqual, "to be deleted")
			})
		})
	})
	Convey("list the team members", t, func() {
		pid := 54
		nsname := "test/abc"
		h := th(false, func(v url.Values) (interface{}, error, int) {
			var m Member
			return []Member{m}, nil, 200
		})
		srv, cl := StubHandler(h)
		defer srv.Close()
		_, e := cl.AllTeamMembers(nil, nil, nil)
		Convey("with no given team-id or name it must shout an error", func() {
			So(e, ShouldNotBeNil)
		})
		m, _ := cl.AllTeamMembers(&pid, nil, nil)
		// copy out the values of "h", because it is reused by the next test
		idpath := h.path
		idmeth := h.method
		Convey("but with a teamid it should return a value", func() {
			So(len(m), ShouldEqual, 1)
			So(idmeth, ShouldEqual, "GET")
			So(idpath, ShouldEqual, fmt.Sprintf("/projects/%d/members", pid))
		})
		m2, _ := cl.AllTeamMembers(nil, &nsname, nil)
		path := h.path
		meth := h.method
		Convey("and with a teamname it should return a value too", func() {
			So(len(m2), ShouldEqual, 1)
			So(meth, ShouldEqual, "GET")
			So(path, ShouldEqual, fmt.Sprintf("/projects/%s/members", nsname))
		})
	})

	Convey("get a specific member", t, func() {
		pid := 1
		uid := 2
		name := "lullaby"
		email := "me@you.com"
		username := "lalluby"
		state := MemberActive
		created := time.Date(2010, time.April, 20, 10, 20, 30, 0, time.UTC)

		h := th(false, func(v url.Values) (interface{}, error, int) {
			var m Member
			m.Username = username
			m.Name = name
			m.EMail = email
			m.Id = uid
			m.State = state
			m.Access = Developer
			m.Created = created
			return &m, nil, 200
		})
		srv, cl := StubHandler(h)
		defer srv.Close()
		m, _ := cl.Member(pid, uid)
		Convey("the url should be correct", func() {
			So(h.method, ShouldEqual, "GET")
			So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/members/%d", pid, uid))
			Convey("and the result should contain correct values", func() {
				So(m.Username, ShouldEqual, username)
				So(m.Name, ShouldEqual, name)
				So(m.EMail, ShouldEqual, email)
				So(m.Id, ShouldEqual, uid)
				So(m.State, ShouldEqual, state)
				So(m.Access, ShouldEqual, Developer)
				So(m.Created.String(), ShouldEqual, created.String())
			})
		})
	})
}
