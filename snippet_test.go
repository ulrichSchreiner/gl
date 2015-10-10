package gl

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestSnippets(t *testing.T) {
	Convey("Snippets functions", t, func() {
		Convey("List all snippets", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []Snippet{Snippet{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.AllSnippets("1")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/snippets")
			})
		})
		Convey("Get a single snippet", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &Snippet{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.GetSnippet("1", 2)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/snippets/2")
			})
		})
		Convey("Create a snippet", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &Snippet{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			title, filename, code := "title", "filename", "code"
			cl.CreateSnippet("1", title, filename, code, Public)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, "/projects/1/snippets")
				So(h.get("title"), ShouldEqual, title)
				So(h.get("file_name"), ShouldEqual, filename)
				So(h.get("code"), ShouldEqual, code)
			})
		})
		Convey("Edit a Snippet", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &Snippet{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			title, filename, code := "title", "filename", "code"
			cl.EditSnippet("1", 2, &title, &filename, &code)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, "/projects/1/snippets/2")
				So(h.get("title"), ShouldEqual, title)
				So(h.get("file_name"), ShouldEqual, filename)
				So(h.get("code"), ShouldEqual, code)
			})
		})
		Convey("Delete a snippet", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &Snippet{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.DeleteSnippet("1", 2)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, "/projects/1/snippets/2")
			})
		})
	})
}
