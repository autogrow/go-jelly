package iso8601

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

/********************************** NOTICE *************************************
* we use UTC here because the test could be run on servers in multiple timezones
*******************************************************************************/

func TestParse(t *testing.T) {
	Convey("given a an ISO8601 string", t, func() {
		t8601 := "2017-03-06T05:05:57.306Z"

		Convey("it should convert it to the correct time", func() {
			t, err := Parse(t8601)
			So(err, ShouldBeNil)

			So(t.Second(), ShouldEqual, 57)
			So(t.Minute(), ShouldEqual, 05)
			So(t.Hour(), ShouldEqual, 05)
			So(t.Year(), ShouldEqual, 2017)
			So(t.Month(), ShouldEqual, 03)
			So(t.Day(), ShouldEqual, 06)
		})
	})

	Convey("given a time", t, func() {

		t := time.Now().UTC()

		Convey("it should convert to an ISO8601 time", func() {
			t8601 := Convert(t)
			exp := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
			So(t8601, ShouldEqual, exp)
		})

	})

	Convey("generate, parse and convert back", t, func() {
		t8601 := Now()
		t, err := Parse(t8601)
		So(err, ShouldBeNil)

		expT8601 := Convert(t)

		So(t8601, ShouldEqual, expT8601)
	})
}
