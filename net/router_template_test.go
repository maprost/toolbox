package net_test

import (
	"testing"

	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/net"
)

func TestCreateHTMLTemplates(t *testing.T) {
	t.Skip("Need a web path")
	router := net.NewRouter()

	tmp := router.CreateTemplatesFromPath("")
	should.NotBeNil(t, tmp)

	index := tmp.Lookup("index.html")
	should.NotBeNil(t, index)

	header := tmp.Lookup("auth/dashboard.html")
	should.NotBeNil(t, header)
}
