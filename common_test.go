package gl

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestJsonDate(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given a specific date", t, func() {
		t := time.Date(2014, time.August, 10, 10, 10, 10, 10, time.UTC)
		jd := JsonDate{t}

		Convey("The representation should only contain the date", func() {
			txt, err := jd.MarshalText()
			So(err, ShouldBeNil)
			So(string(txt), ShouldEqual, "2014-08-10")

			Convey("And the json representation should be included in quotes", func() {
				txt, err := jd.MarshalJSON()
				So(err, ShouldBeNil)
				So(string(txt), ShouldEqual, "\"2014-08-10\"")
			})
		})
	})
	Convey("Given a plain date string", t, func() {
		t := []byte("2014-08-10")
		var jd JsonDate
		err := jd.UnmarshalText(t)
		Convey("it should decode to the correct date", func() {
			So(err, ShouldBeNil)
			So(jd.Time.Day(), ShouldEqual, 10)
			So(jd.Time.Month(), ShouldEqual, time.August)
			So(jd.Time.Year(), ShouldEqual, 2014)
		})
	})
	Convey("Given a json date string", t, func() {
		t := []byte("\"2014-08-10\"")
		var jd JsonDate
		err := jd.UnmarshalJSON(t)
		Convey("it should decode to the correct date", func() {
			So(err, ShouldBeNil)
			So(jd.Time.Day(), ShouldEqual, 10)
			So(jd.Time.Month(), ShouldEqual, time.August)
			So(jd.Time.Year(), ShouldEqual, 2014)
		})
	})
}
