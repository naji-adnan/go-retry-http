package goretryhttp

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBackoff(t *testing.T) {
	Convey("Given RetryWaitMin:1sec, RetryWaitMax:30sec", t, func() {
		min := 1 * time.Second
		max := 30 * time.Second

		Convey("When on first attempt, backoff for 2 sec", func() {
			durationFirst := DefaultBackoff(min, max, 1)
			So(durationFirst, ShouldEqual, 2*time.Second)
		})

		Convey("When on second attempt, backoff for 4 sec", func() {
			durationSecond := DefaultBackoff(min, max, 2)
			So(durationSecond, ShouldEqual, 4*time.Second)
		})

		Convey("When on third attempt, backoff for 8 sec", func() {
			durationThird := DefaultBackoff(min, max, 3)
			So(durationThird, ShouldEqual, 8*time.Second)
		})

		Convey("When on fourth attempt, backoff for 16 sec", func() {
			durationFourth := DefaultBackoff(min, max, 4)
			So(durationFourth, ShouldEqual, 16*time.Second)
		})
	})

	Convey("Given RetryWaitMin:1sec, RetryWaitMax:4sec", t, func() {
		min := 1 * time.Second
		max := 4 * time.Second

		Convey("When on first attempt, backoff for 2 sec", func() {
			durationFirst := DefaultBackoff(min, max, 1)
			So(durationFirst, ShouldEqual, 2*time.Second)
		})

		Convey("When on second attempt, backoff for 4 sec", func() {
			durationSecond := DefaultBackoff(min, max, 2)
			So(durationSecond, ShouldEqual, 4*time.Second)
		})

		Convey("When on third attempt, backoff for 4 sec (max limit) only even if sleep exceeds max limit", func() {
			durationThird := DefaultBackoff(min, max, 3)
			So(durationThird, ShouldEqual, 4*time.Second)
		})
	})

}
