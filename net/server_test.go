package net_test

import (
	"net/http/httptest"

	"github.com/maprost/toolbox/net"
)

type TestStruct struct {
	Name string
}

func NewGetServer(actions map[string]net.HandlerFunc) *httptest.Server {
	//gin.SetMode(gin.ReleaseMode)
	server := net.NewServer(net.NewConfig())

	for path, action := range actions {
		server.Get(path, func(con *net.Connection) {
			action(con)
		})
	}
	testServer := httptest.NewServer(server)
	return testServer
}

func NewSimpleGetServer(action net.HandlerFunc) *httptest.Server {
	return NewGetServer(map[string]net.HandlerFunc{"/": action})
}
