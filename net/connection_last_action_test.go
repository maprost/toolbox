package net_test

import (
	"net/http"
	"testing"

	"github.com/maprost/restclient"
	"github.com/maprost/restclient/rctest"
	"github.com/maprost/testbox/must"
	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/net"
)

func Test_SendResponse_200(t *testing.T) {
	server := NewSimpleGetServer(func(con *net.Connection) {
		ts := TestStruct{Name: "Bob"}
		con.SendResponse(ts, nil)
	})
	defer server.Close()

	var resultTs TestStruct
	response := restclient.Get(server.URL + "/").SendAndGetJsonResponse(&resultTs)
	rctest.CheckResult(t, response, rctest.Status200())

	must.NotBeNil(t, resultTs)
	should.BeEqual(t, resultTs.Name, "Bob")
}

func Test_SendResponse_204(t *testing.T) {
	server := NewSimpleGetServer(func(con *net.Connection) {
		con.SendResponse(nil, nil)
	})
	defer server.Close()

	response := restclient.Get(server.URL + "/").Send()
	rctest.CheckResult(t, response, rctest.Status204())
}

func Test_SendResponse_400(t *testing.T) {
	server := NewSimpleGetServer(func(con *net.Connection) {
		ts := TestStruct{Name: "Bob"}
		con.SendResponse(ts, net.NewBadRequestError("all is wrong here", "Msg"))
	})
	defer server.Close()

	var resultTs TestStruct
	response := restclient.Get(server.URL + "/").SendAndGetJsonResponse(&resultTs)
	rctest.CheckResult(t, response, rctest.FailedResponse(http.StatusBadRequest, `{"brokenElement":"Msg","brokenIndex":0,"msg":"all is wrong here"}`))

	must.NotBeNil(t, resultTs)
	should.BeEqual(t, resultTs.Name, "")
}

func Test_SendResponse_404(t *testing.T) {
	server := NewSimpleGetServer(func(con *net.Connection) {
		ts := TestStruct{Name: "Bob"}
		con.SendResponse(ts, net.NewNotFoundError("all is wrong here"))
	})
	defer server.Close()

	var resultTs TestStruct
	response := restclient.Get(server.URL + "/").SendAndGetJsonResponse(&resultTs)
	rctest.CheckResult(t, response, rctest.FailedResponse(http.StatusNotFound, "error occurred\n"))

	must.NotBeNil(t, resultTs)
	should.BeEqual(t, resultTs.Name, "")
}

func Test_Redirect(t *testing.T) {
	server := NewGetServer(
		map[string]net.HandlerFunc{
			"/redirect": func(con *net.Connection) {
				con.Redirect("/target")
			},
			"/target": func(con *net.Connection) {
				con.SendResponse(nil, nil)
			},
		})
	defer server.Close()

	response := restclient.Get(server.URL + "/redirect").Send()
	rctest.CheckResult(t, response, rctest.Status204())
}

func Test_FailRedirect(t *testing.T) {
	server := NewGetServer(
		map[string]net.HandlerFunc{
			"/redirect": func(con *net.Connection) {
				con.FailRedirect()
			},
			"/": func(con *net.Connection) {
				con.SendResponse(nil, nil)
			},
		})
	defer server.Close()

	response := restclient.Get(server.URL + "/redirect").Send()
	rctest.CheckResult(t, response, rctest.Status204())
}
