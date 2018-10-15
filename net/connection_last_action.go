package net

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Redirect redirects the website to another CookiePath
func (con *Connection) Redirect(path string) {
	con.saveDefaultCookies()

	con.gin.Redirect(http.StatusFound, path)
	con.setResponseInfo(RedirectType, path)
}

// FailRedirect redirects the website to the root CookiePath and delete the session
func (con *Connection) FailRedirect() {
	con.ClearDefaultCookies()

	con.gin.Redirect(http.StatusFound, con.cfg.FailRedirectPath)
	con.setResponseInfo(FailRedirectType, con.cfg.FailRedirectPath)
}

// SendResponse get an returnDto and an error and send the responseComment
// -> error -> send the error
// -> returnDto -> send the dto with 200
// -> both nil -> send a 204
func (con *Connection) SendResponse(returnDto interface{}, err error) {
	// error?
	if err != nil {
		netErr, ok := err.(netError)
		if !ok {
			// fallback
			netErr = NewInternalServerError(err).(netError)
		}

		if netErr.Code() == http.StatusBadRequest {
			con.gin.JSON(netErr.Code(), gin.H{
				"msg":           netErr.Error(),
				"brokenElement": netErr.BrokenElement(),
				"brokenIndex":   netErr.BrokenIndex(),
			})
		} else {
			http.Error(con.gin.Writer, "error occurred", netErr.Code())
		}

		con.setResponseInfo(RestType, netErr.Error()+"\n"+netErr.Stacktrace())
		return
	}

	// 204?
	if returnDto == nil {
		con.gin.Writer.WriteHeader(http.StatusNoContent)
		con.setResponseInfo(RestType, "")
		return
	}

	// 200
	con.gin.JSON(http.StatusOK, returnDto)
	con.setResponseInfo(RestType, fmt.Sprintf("%+v", returnDto))
}

// SendFile load the file from the CookiePath and return it
func (con *Connection) SendFile(path string) {
	con.gin.File(path)
	con.setResponseInfo(FileType, path)
}

// SendHTML return the html file (out of the cache)
func (con *Connection) SendHTML(keys string, data interface{}) {
	con.saveDefaultCookies()

	con.gin.HTML(http.StatusOK, keys, data)
	con.setResponseInfo(HtmlType, keys)
}

// SendFreshHTML load the html file and return it
func (con *Connection) SendFreshHTML(key string, content string, data interface{}) {
	con.saveDefaultCookies()

	tmpl, err := loadTemplate(key, content, con.cfg)
	if err != nil {
		con.SendResponse(nil, NewNotFoundError("Page-Path: %s\nLoad template error: %s", key, err.Error()))
		return
	}

	con.gin.Writer.Header().Add("Content Type", "text/html")
	err = tmpl.ExecuteTemplate(con.gin.Writer, filepath.Base(key), data)

	if err != nil {
		con.SendResponse(nil, NewNotFoundError("Page-Path: %s\nLoad template error: %s", key, err.Error()))
		return
	}

	con.setResponseInfo(FreshHtmlType, key)
}
