package time

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMidday(t *testing.T) {

	Convey("it should get midday", t, func() {
		Convey("when given a timestamp at midnight", func() {
			n := time.Now()
			t := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)

			mepoch := Midday(float64(t.Unix()))
			m := time.Unix(int64(mepoch), 0)
			So(m.Hour(), ShouldEqual, 12)
			So(m.Minute(), ShouldEqual, 0)
		})

		Convey("when given a random timestamp", func() {
			n := time.Now()
			t := time.Date(n.Year(), n.Month(), n.Day(), 13, 44, 0, 0, time.Local)

			mepoch := Midday(float64(t.Unix()))
			m := time.Unix(int64(mepoch), 0)
			So(m.Hour(), ShouldEqual, 12)
			So(m.Minute(), ShouldEqual, 0)
		})
	})

}
