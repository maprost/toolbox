package net

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type responseType int

const (
	RestType = responseType(iota)
	HtmlType
	FreshHtmlType
	FileType
	RedirectType
	FailRedirectType
)

func (r responseType) String() string {
	switch r {
	case RestType:
		return "Rest"
	case HtmlType:
		return "Html"
	case FreshHtmlType:
		return "fresh-Html"
	case FileType:
		return "File"
	case RedirectType:
		return "Redirect"
	case FailRedirectType:
		return "FailRedirect"
	}
	return "unknown"
}

type ResponseInfo struct {
	Type        responseType
	Code        int
	Description string
	Duration    time.Duration
}

// 200,204 or a redirect that was not a failed redirect.
func (r ResponseInfo) Successful() bool {
	return r.Code <= 204 || (r.Code < 400 && r.Type != FailRedirectType)
}

// Connection wraps the gin.Context struct
type Connection struct {
	cfg          *Config
	gin          *gin.Context
	responseInfo ResponseInfo
	start        time.Time
	wsUpgrader   websocket.Upgrader

	defaultCookieValues map[string]string // cookie -> value

	Context interface{}
}

// RequestSignature returns the request signature as string
func (con *Connection) RequestSignature() string {
	return con.gin.Request.Method + ":" + con.gin.Request.RequestURI
}

// RequestHeader is more for debugging and returns the request header as string
func (con *Connection) RequestHeader() string {
	return fmt.Sprint(con.gin.Request.Header)
}

// ResponseHeader is more for debugging and returns the responseComment header as string
func (con *Connection) ResponseHeader() string {
	return fmt.Sprint(con.gin.Writer.Header())
}

// ResponseInfo returns a short notice of the current response
func (con *Connection) ResponseInfo() ResponseInfo {
	return con.responseInfo
}

// setResponseInfo set all attributes of the responseInfo plus the time the request needed
func (con *Connection) setResponseInfo(typ responseType, description string) {
	con.responseInfo = ResponseInfo{
		Type:        typ,
		Code:        con.gin.Writer.Status(),
		Description: description,
		Duration:    time.Since(con.start),
	}
}
