package goretryhttp

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBackoff(t *testing.T) {
	Convey("Given RetryWaitMin:1sec, RetryWaitMax:30sec", t, func() {
		mockResp := http.Response{StatusCode: 200}
		min := 1 * time.Second
		max := 30 * time.Second

		Convey("When on first attempt, backoff for 2 sec", func() {
			durationFirst := DefaultBackoff(min, max, 1, &mockResp)
			So(durationFirst, ShouldEqual, 2*time.Second)
		})

		Convey("When on second attempt, backoff for 4 sec", func() {
			durationSecond := DefaultBackoff(min, max, 2, &mockResp)
			So(durationSecond, ShouldEqual, 4*time.Second)
		})

		Convey("When on third attempt, backoff for 8 sec", func() {
			durationThird := DefaultBackoff(min, max, 3, &mockResp)
			So(durationThird, ShouldEqual, 8*time.Second)
		})

		Convey("When on fourth attempt, backoff for 16 sec", func() {
			durationFourth := DefaultBackoff(min, max, 4, &mockResp)
			So(durationFourth, ShouldEqual, 16*time.Second)
		})

	})
}
