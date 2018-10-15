package net_test

import (
	"testing"

	"github.com/maprost/restclient"
	"github.com/maprost/restclient/rctest"
	"github.com/maprost/testbox/must"
	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/net"
)

func Test_QueryParamDate(t *testing.T) {
	server := NewSimpleGetServer(func(con *net.Connection) {
		from, err := con.QueryParamDate("from", "2006-01-02")
		con.SendResponse(from.String(), err)
	})
	defer server.Close()

	response := restclient.Get(server.URL+"/").AddQueryParam("from", "2018-12-24").SendAndGetResponseItem()
	should.NotBeNil(t, response)
	rctest.CheckResult(t, response.Result, rctest.Status200())

	should.BeEqual(t, response.String(), "\"2018-12-24 00:00:00 +0000 UTC\"")
}

func Test_QueryParamDate_wrongFormat(t *testing.T) {
	server := NewSimpleGetServer(func(con *net.Connection) {
		from, err := con.QueryParamDate("from", "2006-01-02")
		con.SendResponse(from.String(), err)
	})
	defer server.Close()

	response := restclient.Get(server.URL+"/").AddQueryParam("from", "2018 12 24").SendAndGetResponseItem()
	must.NotBeNil(t, response)
	rctest.CheckResult(t, response.Result, rctest.Status404())
}
