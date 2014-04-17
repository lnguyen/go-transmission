package client

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux *http.ServeMux

	client ApiClient

	server *httptest.Server
)

func RPCHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Transmission-Session-Id") == "" {
		res.Header().Set("X-Transmission-Session-Id", "123")
		res.WriteHeader(http.StatusConflict)
		return
	}
	fmt.Fprintf(res, `{"arguments":{},"result":"no method name"}`)
}

func setup() {
	// test server
	mux = http.NewServeMux()
	m := martini.New()
	r := martini.NewRouter()
	r.Post("/transmission/rpc", RPCHandler)
	m.Action(r.Handle)
	m.Use(auth.Basic("test", "test"))
	mux.Handle("/", m)
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient(server.URL, "test", "test")
}

func teardown() {
	server.Close()
}

func TestPost(t *testing.T) {
	setup()
	defer teardown()
	Convey("Test Post is working correctly", t, func() {
		output, err := client.Post("")
		So(err, ShouldBeNil)
		So(string(output), ShouldEqual, `{"arguments":{},"result":"no method name"}`)
	})

	Convey("Test when auth is incorrect", t, func() {
		fakeClient := NewClient(server.URL, "testfake", "testfake")
		output, err := fakeClient.Post("")
		So(err, ShouldBeNil)
		So(string(output), ShouldEqual, "Not Authorized\n")
	})

}
