package net_test

import (
	"net/http/httptest"

)

type actionFunc func(net.Connection)

type TestStruct struct {
	Name string
}

func NewNetServer(actions map[string]actionFunc) *httptest.Server {
	//gin.SetMode(gin.ReleaseMode)
	router := net.NewRouter()

	for path, action := range actions {
		router.Get(path, func(net net.Connection) {
			action(net)
		})
	}
	server := httptest.NewServer(router)
	return server
}

func NewSimpleNetServer(action actionFunc) *httptest.Server {
	return NewNetServer(map[string]actionFunc{"/": action})
}
