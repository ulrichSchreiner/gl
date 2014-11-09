package gl

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestRepository(t *testing.T) {
	Convey("repository test functions", t, func() {
		Convey("List repository tags", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []TagListEntry{TagListEntry{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.AllTags("1")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/tags")
			})
		})
		Convey("create a repository tag", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &Tag{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			msg := "tag message"
			cl.CreateTag("1", "tagname", "tag ref", &msg)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, "/projects/1/repository/tags")
				So(h.get("tag_name"), ShouldEqual, "tagname")
				So(h.get("ref"), ShouldEqual, "tag ref")
				So(h.get("message"), ShouldEqual, msg)
			})
		})
		Convey("List repository entries", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []RepositoryEntry{RepositoryEntry{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			pt := "repopath"
			ref := "reporef"
			cl.AllRepoEntries("1", &pt, &ref)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/tree")
				So(h.get("path"), ShouldEqual, pt)
				So(h.get("ref_name"), ShouldEqual, ref)
			})
		})
		Convey("get the raw content of a file", func() {
			h := thp(func(v url.Values) (interface{}, error, int) {
				return []byte("hello"), nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			res, _ := cl.RawFileContent("1", "sha1", "apath")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/blobs/sha1")
				So(h.get("filepath"), ShouldEqual, "apath")
				So(string(res), ShouldEqual, "hello")
			})
		})
		Convey("get the raw blob content of a file", func() {
			h := thp(func(v url.Values) (interface{}, error, int) {
				return []byte("hello"), nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			res, _ := cl.RawBlobContent("1", "sha1")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/raw_blobs/sha1")
				So(string(res), ShouldEqual, "hello")
			})
		})
		Convey("get an archive", func() {
			h := thp(func(v url.Values) (interface{}, error, int) {
				return []byte("hello"), nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			sha := "sha1"
			res, _ := cl.Archive("1", &sha)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/archive")
				So(h.get("sha"), ShouldEqual, sha)
				So(string(res), ShouldEqual, "hello")
			})
		})
		Convey("compare two revisions", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &Comparison{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.Compare("1", "start", "end")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/compare")
				So(h.get("from"), ShouldEqual, "start")
				So(h.get("to"), ShouldEqual, "end")
			})
		})
		Convey("read all contributors", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []Contributor{Contributor{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.AllContributors("1")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/contributors")
			})
		})
	})
	Convey("repository crud for files", t, func() {
		Convey("read a specific file", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return RepoFile{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.ReadFile("1", "myfile", "myref")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/files")
				So(h.get("file_path"), ShouldEqual, "myfile")
				So(h.get("ref"), ShouldEqual, "myref")
			})
		})
		Convey("create a file", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return RepoFile{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.CreateFile("1", "myfile", "mybranch", "mycommit", "mycontent", "myenc")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, "/projects/1/repository/files")
				So(h.get("file_path"), ShouldEqual, "myfile")
				So(h.get("branch_name"), ShouldEqual, "mybranch")
				So(h.get("commit_message"), ShouldEqual, "mycommit")
				So(h.get("content"), ShouldEqual, "mycontent")
				So(h.get("encoding"), ShouldEqual, "myenc")
			})
		})
		Convey("update a file", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return RepoFile{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.UpdateFile("1", "myfile", "mybranch", "mycommit", "mycontent", "myenc")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, "/projects/1/repository/files")
				So(h.get("file_path"), ShouldEqual, "myfile")
				So(h.get("branch_name"), ShouldEqual, "mybranch")
				So(h.get("commit_message"), ShouldEqual, "mycommit")
				So(h.get("content"), ShouldEqual, "mycontent")
				So(h.get("encoding"), ShouldEqual, "myenc")
			})
		})
		Convey("delete a file", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return RepoFile{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.DeleteFile("1", "myfile", "mybranch", "mycommit")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, "/projects/1/repository/files")
				So(h.get("file_path"), ShouldEqual, "myfile")
				So(h.get("branch_name"), ShouldEqual, "mybranch")
				So(h.get("commit_message"), ShouldEqual, "mycommit")
			})
		})
	})
	Convey("check the commit api", t, func() {
		Convey("receive all commits", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []Commit{Commit{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			ref := "mybranch"
			cl.AllCommits("1", &ref)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/commits")
				So(h.get("ref_name"), ShouldEqual, "mybranch")
			})
		})
		Convey("read one commit", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return Commit{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.ReadCommit("1", "mysha")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/commits/mysha")
			})
		})
		Convey("read a diff", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return Diff{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.ReadDiff("1", "mysha")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/repository/commits/mysha/diff")
			})
		})
	})
}
