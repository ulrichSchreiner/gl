package gl

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestProjectCreation(t *testing.T) {
	Convey("Creating a project with only a name", t, func() {
		name := "testproject"

		var vals url.Values
		srv, cl := StubHandler(t, "POST", true, func(t *testing.T, v url.Values) (interface{}, error, int) {
			var res Project
			res.Name = name
			vals = v
			return &res, nil, 200
		})
		defer srv.Close()
		p, _ := cl.CreateProject(name, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		Convey("the mapped URL parameters must be correct", func() {
			So(vals.Get("name"), ShouldEqual, name)
			So(vals, hasnot,
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

		var vals url.Values
		srv, cl := StubHandler(t, "POST", true, func(t *testing.T, v url.Values) (interface{}, error, int) {
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

			vals = v
			return &res, nil, 200
		})
		defer srv.Close()
		p, _ := cl.CreateProject(name, &path, &namespaceid, &desc, &iss, &merge, &wiki, &snipp, &pub, &vis, &iurl)
		Convey("the mapped URL parameters must be correct", func() {
			So(vals.Get("name"), ShouldEqual, name)
			So(vals, has,
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
}
