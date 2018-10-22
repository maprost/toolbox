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

// Request without auth check
func WebSocket(action WebSocketFunc) HandlerFunc {
	return webSocket(noCheck, func(*Connection) {}, action)
}

// Request without auth check
func AuthWebSocket(action WebSocketFunc) HandlerFunc {
	return webSocket(authCheck, failedAuthRequest, action)
}

// Request without auth check
func AdminWebSocket(action WebSocketFunc) HandlerFunc {
	return webSocket(adminCheck, failedAuthRequest, action)
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

func request(checkAction CheckFunc, failAuthAction HandlerFunc, action HandlerFunc) HandlerFunc {
	if checkAction == nil || failAuthAction == nil || action == nil {
		panic(fmt.Sprintf("can't execute request method, some given methods are nil (checkAction: %t, failAuthAction: %t, action: %t)", checkAction == nil, failAuthAction == nil, action == nil))
	}

	return func(con *Connection) {
		if con.cfg.Close == nil || con.cfg.Finish == nil {
			panic(fmt.Sprintf("can't execute request method, some given methods are nil (close: %t, finish: %t)", con.cfg.Close == nil, con.cfg.Finish == nil))
		}

		log.Print("\n\nRequest: ", con.RequestSignature())
		log.Print(con.RequestHeader())

		// security check
		ok := checkAction(con)
		if !ok {
			log.Print("Check failed")
			failAuthAction(con)
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

func webSocket(checkAction CheckFunc, failAuthAction HandlerFunc, action WebSocketFunc) HandlerFunc {
	if checkAction == nil || failAuthAction == nil || action == nil {
		panic(fmt.Sprintf("can't execute request method, some given methods are nil (checkAction: %t, failAuthAction: %t, action: %t)", checkAction == nil, failAuthAction == nil, action == nil))
	}

	return func(con *Connection) {
		if con.cfg.Close == nil || con.cfg.Finish == nil || con.cfg.WebSocketError == nil {
			panic(fmt.Sprintf("can't execute request method, some given methods are nil (close: %t, finish: %t, webSocketError: %t)", con.cfg.Close == nil, con.cfg.Finish == nil, con.cfg.WebSocketError == nil))
		}

		ws, err := con.WebSocketChannel()
		if err != nil {
			con.cfg.WebSocketError(con)
			return
		}

		defer ws.Close()

		// security check
		ok := checkAction(con)
		if !ok {
			log.Print("Check failed")
			failAuthAction(con)
		} else {
			for {

				msg, open, err := ws.Read()
				if err != nil {
					con.cfg.WebSocketError(con)
					return
				}
				if !open {
					con.cfg.Finish(con)
					return
				}

				con.Context = con.cfg.InitContext()

				action(con)

				response := con.ResponseInfo()
				log.Printf("ResponseInfo: %s[%d](%s) %s", response.Type, response.Code, response.Duration, response.Description)

				commitChanges := response.Code <= 204 || (response.Code < 400 && response.Type != FailRedirectType)
				err := con.cfg.Close(con, commitChanges)
				if err != nil {
					log.Print("Context Close error")
				}
			}
		}

		log.Print("\n\nRequest: ", con.RequestSignature())
		log.Print(con.RequestHeader())

		// security check
		ok := checkAction(con)
		if !ok {
			log.Print("Check failed")
			failAuthAction(con)
		} else {
			action(con)
		}

		con.cfg.Finish(con)
	}
}
