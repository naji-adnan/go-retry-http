package goretryhttp

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultRetryPolicy(t *testing.T) {

	Convey("Given no connection errors", t, func() {

		Convey("When the server response StatusCode is Success (2xx StatusCode)", func() {
			mockResp200 := http.Response{StatusCode: 200}

			Convey("retrying request is not activated", func() {
				shouldRetry, err := DefaultRetryPolicy(&mockResp200, nil)

				So(shouldRetry, ShouldEqual, false)
				So(err, ShouldEqual, nil)
			})
		})

		Convey("When the server response StatusCode is 0", func() {
			mockResp0 := http.Response{StatusCode: 0}

			Convey("retrying request is activated", func() {
				shouldRetry, err := DefaultRetryPolicy(&mockResp0, nil)

				So(shouldRetry, ShouldEqual, true)
				So(err, ShouldEqual, nil)
			})
		})

		Convey("When the server response StatusCode is 5xx", func() {
			mockResp0 := http.Response{StatusCode: 500}

			Convey("retrying request is activated", func() {
				shouldRetry, err := DefaultRetryPolicy(&mockResp0, nil)

				So(shouldRetry, ShouldEqual, true)
				So(err, ShouldEqual, nil)
			})
		})

		Convey("When the server response StatusCode is not Success and not 5xx", func() {
			mockResp0 := http.Response{StatusCode: 404}

			Convey("retrying request is not activated", func() {
				shouldRetry, err := DefaultRetryPolicy(&mockResp0, nil)

				So(shouldRetry, ShouldEqual, false)
				So(err, ShouldEqual, nil)
			})
		})
	})

	Convey("Given a connection error", t, func() {

		Convey("When the server response StatusCode is Success (2xx StatusCode)", func() {
			mockResp200 := http.Response{StatusCode: 200}
			mockErr := errors.New("a mock error occurred")

			Convey("retrying request is activated", func() {
				shouldRetry, err := DefaultRetryPolicy(&mockResp200, mockErr)

				So(shouldRetry, ShouldEqual, true)
				So(err, ShouldEqual, mockErr)
			})
		})
	})

}
