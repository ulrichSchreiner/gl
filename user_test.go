package gl

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestUsers(t *testing.T) {
	Convey("Users functions", t, func() {
		Convey("Search for users", func() {
			name := "searchfor"
			h := th(func(v url.Values) (interface{}, error, int) {
				return []User{User{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.SearchAllUsers(name)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, users_url)
				So(h.get("search"), ShouldEqual, name)
			})
		})
		Convey("List all users", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []User{User{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.AllUsers()
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, users_url)
			})
		})
		Convey("Get a single user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &User{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.GetUser(1)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/users/1")
			})
		})
		Convey("Get current user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &User{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.CurrentUser()
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/user")
			})
		})
		Convey("Create a user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &User{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			email, name, pass, username := "email", "name", "password", "username"
			skype, twitter, linkedin, website := "skype", "twitter", "linkedin", "website"
			limit := 10
			externuid, provider, bio := "externuid", "provider", "bio"
			admin, cancreate := true, true

			cl.CreateUser(email, username, pass, name, &skype, &linkedin, &twitter, &website, &limit, &externuid, &provider, &bio, &admin, &cancreate)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, "/users")
				So(h.get("email"), ShouldEqual, email)
				So(h.get("name"), ShouldEqual, name)
				So(h.get("password"), ShouldEqual, pass)
				So(h.get("username"), ShouldEqual, username)
				So(h.get("skype"), ShouldEqual, skype)
				So(h.get("twitter"), ShouldEqual, twitter)
				So(h.get("linkedin"), ShouldEqual, linkedin)
				So(h.get("website_url"), ShouldEqual, website)
				So(h.get("projects_limit"), ShouldEqual, fmt.Sprintf("%v", limit))
				So(h.get("can_create_group"), ShouldEqual, fmt.Sprintf("%v", cancreate))
				So(h.get("admin"), ShouldEqual, fmt.Sprintf("%v", admin))
				So(h.get("extern_uid"), ShouldEqual, externuid)
				So(h.get("provider"), ShouldEqual, provider)
				So(h.get("bio"), ShouldEqual, bio)
			})
		})
		Convey("Edit a user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &User{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			email, name, pass, username := "email", "name", "password", "username"
			skype, twitter, linkedin, website := "skype", "twitter", "linkedin", "website"
			limit := 10
			externuid, provider, bio := "externuid", "provider", "bio"
			admin, cancreate := true, true

			cl.EditUser(4, email, username, pass, name, &skype, &linkedin, &twitter, &website, &limit, &externuid, &provider, &bio, &admin, &cancreate)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "PUT")
				So(h.path, ShouldEqual, "/users/4")
				So(h.get("email"), ShouldEqual, email)
				So(h.get("name"), ShouldEqual, name)
				So(h.get("password"), ShouldEqual, pass)
				So(h.get("username"), ShouldEqual, username)
				So(h.get("skype"), ShouldEqual, skype)
				So(h.get("twitter"), ShouldEqual, twitter)
				So(h.get("linkedin"), ShouldEqual, linkedin)
				So(h.get("website_url"), ShouldEqual, website)
				So(h.get("projects_limit"), ShouldEqual, fmt.Sprintf("%v", limit))
				So(h.get("can_create_group"), ShouldEqual, fmt.Sprintf("%v", cancreate))
				So(h.get("admin"), ShouldEqual, fmt.Sprintf("%v", admin))
				So(h.get("extern_uid"), ShouldEqual, externuid)
				So(h.get("provider"), ShouldEqual, provider)
				So(h.get("bio"), ShouldEqual, bio)
			})
		})
		Convey("Delete a user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &User{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.DeleteUser(4)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, "/users/4")
			})
		})
		Convey("Create a session", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &User{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.Session("login", "email", "pass")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, "/session")
				So(h.get("login"), ShouldEqual, "login")
				So(h.get("email"), ShouldEqual, "email")
				So(h.get("password"), ShouldEqual, "pass")
			})
		})
	})
	Convey("sshkey functions", t, func() {
		Convey("list ssh keys of current user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []SshKey{SshKey{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.AllCurrentUserKeys()
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/user/keys")
			})
		})
		Convey("get a key of current user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return []SshKey{SshKey{}}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.GetSshKey(3)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "GET")
				So(h.path, ShouldEqual, "/user/keys/3")
			})
		})
		Convey("add a key of current user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &SshKey{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.CreateCurrentUserSshKey("keytitle", "keykey")
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "POST")
				So(h.path, ShouldEqual, "/user/keys")
				So(h.get("title"), ShouldEqual, "keytitle")
				So(h.get("key"), ShouldEqual, "keykey")
			})
		})
		Convey("del a key of current user", func() {
			h := th(func(v url.Values) (interface{}, error, int) {
				return &SshKey{}, nil, 200
			})
			srv, cl := StubHandler(h)
			defer srv.Close()
			cl.DeleteCurrentUserKey(1)
			Convey("check if the request was correct", func() {
				So(h.method, ShouldEqual, "DELETE")
				So(h.path, ShouldEqual, "/user/keys/1")
			})
		})
	})
}
