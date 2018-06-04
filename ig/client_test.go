package ig

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIGClient(t *testing.T) {
	Convey("get api tokens", t, func() {

		username := os.Getenv("IG_USERNAME")
		password := os.Getenv("IG_PASSWORD")

		c, err := NewClient(username, password)

		Convey("new client shouldn't be empty", func() {
			So(err, ShouldBeNil)
			So(c.token, ShouldNotBeNil)
			So(c.refreshTime, ShouldNotEqual, 0)
			So(c.refreshToken, ShouldNotBeNil)

			Convey("get devices", func() {
				err := c.GetDevices()
				So(err, ShouldBeNil)

				Convey("get device info", func() {
					c.UpdateAllGrowrooms()
					So(err, ShouldBeNil)
					Convey("test get climate", func() {
						err = c.UpdateGrowroom("1")
						So(err, ShouldBeNil)
						genGet, readErr := c.GetGrowroomReading("1", grAirTemp)
						So(genGet, ShouldNotBeEmpty)
						So(readErr, ShouldBeNil)
						growroom, grErr := c.GetGrowroom("1")
						So(growroom, ShouldNotBeNil)
						So(grErr, ShouldBeNil)
						valid, specGet := growroom.AirTemp()
						So(valid, ShouldBeTrue)
						So(specGet, ShouldNotBeEmpty)
						So(genGet, ShouldEqual, specGet)
					})
					Convey("test get doser", func() {
						err = c.UpdateGrowroom("1")
						So(err, ShouldBeNil)
						genGet, readErr := c.GetGrowroomReading("1", grEC)
						So(genGet, ShouldNotBeEmpty)
						So(readErr, ShouldBeNil)
						growroom, grErr := c.GetGrowroom("1")
						So(growroom, ShouldNotBeNil)
						So(grErr, ShouldBeNil)
						valid, specGet := growroom.EC()
						So(valid, ShouldBeTrue)
						So(specGet, ShouldNotBeEmpty)
						So(genGet, ShouldEqual, specGet)
					})
					Convey("test update doser setting", func() {
						for _, id := range c.MyDevices.Dosers() {
							err = id.GetConfigState()
							So(err, ShouldBeNil)

							dt := id.Status.General.NutrientDoseTime

							newDT := dt + 10
							id.Status.General.NutrientDoseTime = newDT

							err = id.UpdateState()
							So(err, ShouldBeNil)

							dt = id.Status.General.NutrientDoseTime
							So(int(dt), ShouldEqual, newDT)
						}
					})
				})
			})
		})
	})
}
