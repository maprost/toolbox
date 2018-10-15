package net

import (
	"fmt"
	"log"
)

// Request without auth check
func Request(action HandlerFunc) HandlerFunc {
	return request(noCheck, func(*Connection) {}, action)
}

// AuthRequest with auth check and 401 if auth failed
func AuthRequest(action HandlerFunc) HandlerFunc {
	return request(authCheck, failedAuthRequest, action)
}

// AdminRequest with admin check and 401 if auth failed
func AdminRequest(action HandlerFunc) HandlerFunc {
	return request(adminCheck, failedAuthRequest, action)
}

// Website without auth check
func Website(action HandlerFunc) HandlerFunc {
	return request(noCheck, func(*Connection) {}, action)
}

// Website with auth check and root redirect if auth failed
func AuthWebsite(action HandlerFunc) HandlerFunc {
	return request(authCheck, failedAuthWebsite, action)
}

// Website with admin check and root redirect if auth failed
func AdminWebsite(action HandlerFunc) HandlerFunc {
	return request(adminCheck, failedAuthWebsite, action)
}

// File without auth check, no locale check, and no auth failed function
func File(action HandlerFunc) HandlerFunc {
	return request(noCheck, func(*Connection) {}, action)
}

func noCheck(*Connection) bool {
	return true
}

func authCheck(con *Connection) bool {
	return con.cfg.AuthCheck(con)
}

func adminCheck(con *Connection) bool {
	return con.cfg.AdminCheck(con)
}

func failedAuthRequest(con *Connection) {
	con.SendResponse(nil, NewUnauthorizedError("session is expired"))
}

func failedAuthWebsite(con *Connection) {
	con.FailRedirect()
}

func request(checkAction CheckFunc, failAction HandlerFunc, action HandlerFunc) HandlerFunc {
	return func(con *Connection) {
		if checkAction == nil || failAction == nil || action == nil {
			panic(fmt.Sprintf("can't execute request method, some given methods are nil (checkAction: %t, failAction: %t, action: %t)", checkAction == nil, failAction == nil, action == nil))
		}
		if con.cfg.Close == nil || con.cfg.Finish == nil {
			panic(fmt.Sprintf("can't execute request method, some given methods are nil (close: %t, finish: %t)", con.cfg.Close == nil, con.cfg.Finish == nil))
		}

		log.Print("\n\nRequest: ", con.RequestSignature())
		log.Print(con.RequestHeader())

		// security check
		ok := checkAction(con)
		if !ok {
			log.Print("Check failed")
			failAction(con)
		} else {
			action(con)
		}

		response := con.ResponseInfo()
		log.Printf("ResponseInfo: %s[%d](%s) %s", response.Type, response.Code, response.Duration, response.Description)

		commitChanges := response.Code <= 204 || (response.Code < 400 && response.Type != FailRedirectType)
		err := con.cfg.Close(con, commitChanges)
		if err != nil {
			log.Print("Context Close error")
		}

		con.cfg.Finish(con)
	}
}
