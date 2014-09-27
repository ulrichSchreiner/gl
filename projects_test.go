package gl

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestProjects(t *testing.T) {
	Convey("Project functions", t, func() {
		var prj Project
		prj.Name = "testproject"
		prj.Path = "p"
		prj.Description = "desc"
		prj.Public = true
		prj.MergeRequestsEnabled = false
		prj.IssuesEnabled = true
		prj.WikiEnabled = false
		prj.SnippetsEnabled = true
		prj.Visibility = Internal

		Convey("Creating a project with only a name", func() {
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
		Convey("Creating a full filled project", func() {
			namespaceid := 1
			iurl := "import-url"

			h := th(true, func(v url.Values) (interface{}, error, int) {
				return &Project{}, nil, 200
			})

			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.CreateProject(prj.Name, &prj.Path, &namespaceid, &prj.Description,
				&prj.IssuesEnabled, &prj.MergeRequestsEnabled, &prj.WikiEnabled, &prj.SnippetsEnabled, &prj.Public, &prj.Visibility, &iurl)
			Convey("the mapped URL parameters must be correct", func() {
				So(h.get("name"), ShouldEqual, prj.Name)
				So(h.get("path"), ShouldEqual, prj.Path)
				So(h.get("namespace_id"), ShouldEqual, fmt.Sprintf("%d", namespaceid))
				So(h.get("description"), ShouldEqual, prj.Description)
				So(h.get("merge_requests_enabled"), ShouldEqual, fmt.Sprintf("%v", prj.MergeRequestsEnabled))
				So(h.get("issues_enabled"), ShouldEqual, fmt.Sprintf("%v", prj.IssuesEnabled))
				So(h.get("wiki_enabled"), ShouldEqual, fmt.Sprintf("%v", prj.WikiEnabled))
				So(h.get("snippets_enabled"), ShouldEqual, fmt.Sprintf("%v", prj.SnippetsEnabled))
				So(h.get("public"), ShouldEqual, fmt.Sprintf("%v", prj.Public))
				So(h.get("visibility_level"), ShouldEqual, fmt.Sprintf("%v", prj.Visibility))
				So(h.get("import_url"), ShouldEqual, iurl)
			})
		})

		Convey("Search for projects", func() {
			name := "searchfor"
			h := th(false, func(v url.Values) (interface{}, error, int) {
				return []Project{Project{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.SearchAll(name)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/search/%s", name))
			})
		})
		Convey("Creating a userproject with a name and a defaultbranch", func() {
			name := "testproject"
			user := 5
			branch := "testbranch"
			h := th(true, func(v url.Values) (interface{}, error, int) {
				return &Project{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.CreateUserProject(name, user, nil, &branch,
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
		})

		Convey("removing a given project", func() {
			pid := 54
			h := th(true, func(v url.Values) (interface{}, error, int) {
				var p Project
				p.Name = "to be deleted"
				return &p, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.RemoveProject(pid)
			Convey("should invoke a DELETE on the correct url", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, "/projects/54")
			})
		})
	})
	Convey("test the team members functions", t, func() {
		Convey("list the team members", func() {
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
			cl.AllTeamMembers(&pid, nil, nil)
			// copy out the values of "h", because it is reused by the next test
			idpath := h.path
			idmeth := h.method
			Convey("but with a teamid it should return a value", func() {
				So(idmeth, ShouldEqual, "GET")
				So(idpath, ShouldEqual, fmt.Sprintf("/projects/%d/members", pid))
			})
			cl.AllTeamMembers(nil, &nsname, nil)
			path := h.path
			meth := h.method
			Convey("and with a teamname it should return a value too", func() {
				So(meth, ShouldEqual, "GET")
				So(path, ShouldEqual, fmt.Sprintf("/projects/%s/members", nsname))
			})
		})

		Convey("get a specific member", func() {
			pid := 1
			uid := 2

			h := th(false, func(v url.Values) (interface{}, error, int) {
				var m Member
				return &m, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.TeamMember(pid, uid)
			Convey("the url should be correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/members/%d", pid, uid))
			})
		})
		Convey("create a member", func() {
			pid := 54
			h := th(true, func(v url.Values) (interface{}, error, int) {
				var m Member
				return &m, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			_, e := cl.AddTeamMember(nil, nil, 1, Developer)
			Convey("with no given team-id or name it must shout an error", func() {
				So(e, ShouldNotBeNil)
			})
			cl.AddTeamMember(&pid, nil, 12, Master)
			Convey("the url, parameters and return values should be ok", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/members", pid))
				So(h.get("access_level"), ShouldEqual, fmt.Sprintf("%d", Master))
				So(h.get("user_id"), ShouldEqual, "12")
			})
		})
		Convey("edit a member", func() {
			pid := 54
			memb := 12
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var m Member
				return &m, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			_, e := cl.EditTeamMember(nil, nil, 1, Developer)
			Convey("with no given team-id or name it must shout an error", func() {
				So(e, ShouldNotBeNil)
			})
			cl.EditTeamMember(&pid, nil, memb, Master)
			Convey("the url, parameters and return values should be ok", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/members/%d", pid, memb))
				So(h.get("access_level"), ShouldEqual, fmt.Sprintf("%d", Master))
			})
		})
		Convey("delete a member", func() {
			pid := 54
			memb := 12
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var m Member
				return &m, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			_, e := cl.DeleteTeamMember(nil, nil, 1)
			Convey("with no given projectid or name it must shout an error", func() {
				So(e, ShouldNotBeNil)
			})
			cl.DeleteTeamMember(&pid, nil, memb)
			Convey("the url and return values should be ok", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/members/%d", pid, memb))
			})
		})
	})
	Convey("now the project hooks", t, func() {
		Convey("list the team hooks", func() {
			pid := 54
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var h Hook
				return []Hook{h}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			_, e := cl.AllHooks(nil, nil)
			Convey("with no given team-id or name it must shout an error", func() {
				So(e, ShouldNotBeNil)
			})
			cl.AllHooks(&pid, nil)
			Convey("but with a teamid it should return a value", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/hooks", pid))
			})
		})
		Convey("get a specific hook", func() {
			pid := 1
			hid := 2
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var h Hook
				return &h, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.Hook(&pid, nil, hid)
			Convey("the url should be correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/hooks/%d", pid, hid))
			})
		})
		Convey("create a hook", func() {
			pid := 54
			hurl := "myhookurl"
			push, iss, merge := true, false, true
			h := th(true, func(v url.Values) (interface{}, error, int) {
				var h Hook
				return &h, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			_, e := cl.AddHook(nil, nil, "", nil, nil, nil)
			Convey("with no given projectid or name it must shout an error", func() {
				So(e, ShouldNotBeNil)
			})
			cl.AddHook(&pid, nil, hurl, &push, &iss, &merge)
			Convey("the url, parameters and return values should be ok", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/hooks", pid))
				So(h.get("url"), ShouldEqual, hurl)
				So(h.get("push_events"), ShouldEqual, fmt.Sprintf("%v", push))
				So(h.get("issues_events"), ShouldEqual, fmt.Sprintf("%v", iss))
				So(h.get("merge_requests_events"), ShouldEqual, fmt.Sprintf("%v", merge))
			})
		})
		Convey("edit a hook", func() {
			pid := 54
			hid := 65
			hurl := "myhookurl"
			push, iss, merge := true, false, true
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var h Hook
				return &h, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			_, e := cl.EditHook(nil, nil, hid, hurl, nil, nil, nil)
			Convey("with no given projectid or name it must shout an error", func() {
				So(e, ShouldNotBeNil)
			})
			cl.EditHook(&pid, nil, hid, hurl, &push, &iss, &merge)
			Convey("the url, parameters and return values should be ok", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/hooks/%d", pid, hid))
				So(h.get("url"), ShouldEqual, hurl)
				So(h.get("push_events"), ShouldEqual, fmt.Sprintf("%v", push))
				So(h.get("issues_events"), ShouldEqual, fmt.Sprintf("%v", iss))
				So(h.get("merge_requests_events"), ShouldEqual, fmt.Sprintf("%v", merge))
			})
		})
		Convey("delete a hook", func() {
			pid := 54
			hid := 12
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var h Hook
				return &h, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.DeleteHook(&pid, nil, hid)
			Convey("the url and return values should be ok", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/hooks/%d", pid, hid))
			})
		})
	})
	Convey("test the branch services", t, func() {
		Convey("query all branches", func() {
			pid := 54
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var b Branch
				return []Branch{b}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.AllBranches(&pid, nil)
			Convey("check the values of the http request", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/repository/branches", pid))
			})
		})
		Convey("query one specific branch", func() {
			pid := 54
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var b Branch
				return &b, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.Branch(&pid, nil, "bname")
			Convey("check the values of the http request", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/repository/branches/bname", pid))
			})
		})
		Convey("protect a branch", func() {
			pid := 54
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var b Branch
				return &b, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.ProtectBranch(&pid, nil, "bname")
			Convey("check the values of the http request", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/repository/branches/bname/protect", pid))
			})
		})
		Convey("unprotect a branch", func() {
			pid := 54
			h := th(false, func(v url.Values) (interface{}, error, int) {
				var b Branch
				return &b, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.UnprotectBranch(&pid, nil, "bname")
			Convey("check the values of the http request", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/repository/branches/bname/unprotect", pid))
			})
		})
	})
	Convey("now create and delete a fork", t, func() {
		Convey("first create a fork", func() {
			pid := 54
			from := 55
			h := th(false, func(v url.Values) (interface{}, error, int) {
				return nil, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			e := cl.CreateFork(pid, from)
			So(e, ShouldBeNil)
			Convey("check the values of the http request", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/fork/%d", pid, from))
			})
		})
		Convey("delete a fork", func() {
			pid := 54
			h := th(true, func(v url.Values) (interface{}, error, int) {
				return nil, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			e := cl.DeleteFork(pid)
			So(e, ShouldBeNil)
			Convey("check the values of the http request", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, fmt.Sprintf("/projects/%d/fork", pid))
			})
		})
	})
}
