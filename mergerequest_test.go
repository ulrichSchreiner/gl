package gl

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestMergerequests(t *testing.T) {
	Convey("Merge Request functions", t, func() {
		Convey("List all mergerequests", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []MergeRequest{MergeRequest{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			s := OpenedMerges
			ord := OrderByCreated
			cl.AllMergeRequests("1", &s, &ord, true)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/merge_requests")
				So(h.get("state"), ShouldEqual, "opened")
				So(h.get("order_by"), ShouldEqual, "created_at")
				So(h.get("sort"), ShouldEqual, "asc")
			})
		})
		Convey("get a mergerequests", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return MergeRequest{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.GetMergeRequest("1", 2)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/projects/1/merge_request/2")
			})
		})
	})
}
