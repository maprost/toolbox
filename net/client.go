package net

import (
	"net/url"

	"github.com/maprost/restclient"
)

// Client contains all calls toward the rest endpoints of the server
type Client struct {
	BasePath  string
	CookieMap map[string]string
}

func (c Client) addCookies(rc *restclient.RestClient) {
	for cookie, value := range c.CookieMap {
		if value != "" {
			rc.AddHeader("Cookie", cookie+"="+url.QueryEscape(value))
		}
	}
}

func (c Client) Get(path string) *restclient.RestClient {
	rc := restclient.Get(path)
	c.addCookies(rc)
	return rc
}

func (c Client) Post(path string) *restclient.RestClient {
	rc := restclient.Post(path)
	c.addCookies(rc)
	return rc
}

func (c Client) Put(path string) *restclient.RestClient {
	rc := restclient.Put(path)
	c.addCookies(rc)
	return rc
}

func (c Client) Delete(path string) *restclient.RestClient {
	rc := restclient.Delete(path)
	c.addCookies(rc)
	return rc
}
