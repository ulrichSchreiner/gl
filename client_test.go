package gl

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestClientCreation(t *testing.T) {
	Convey("Creating a client with specific settings", t, func() {
		host := "https://myhost"
		api := "/myapi"
		c, err := New(host, api, false)
		Convey("the fields should be initialized correctly", func() {
			So(err, ShouldBeNil)
			So(c.hostURL.String(), ShouldEqual, host)
			So(c.apiPath, ShouldEqual, api)
			t := c.client.Transport.(*http.Transport)
			So(t.TLSClientConfig.InsecureSkipVerify, ShouldBeTrue)
		})
	})
	Convey("Create a secure client", t, func() {
		host := "https://myhost"
		api := "/myapi"
		c, err := Open(host, api)
		Convey("it should be secure", func() {
			So(err, ShouldBeNil)
			So(c.hostURL.String(), ShouldEqual, host)
			So(c.apiPath, ShouldEqual, api)
			t := c.client.Transport.(*http.Transport)
			So(t.TLSClientConfig.InsecureSkipVerify, ShouldBeFalse)
		})
	})
	Convey("Create a v3 client", t, func() {
		host := "https://myhost"
		c, err := OpenV3(host)
		Convey("it should be secure and have a correct api path", func() {
			So(err, ShouldBeNil)
			So(c.hostURL.String(), ShouldEqual, host)
			So(c.apiPath, ShouldEqual, APIv3)
			t := c.client.Transport.(*http.Transport)
			So(t.TLSClientConfig.InsecureSkipVerify, ShouldBeFalse)
		})
	})
}
