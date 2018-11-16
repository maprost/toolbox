package net

import (
	"fmt"
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

// AuthWebsite with auth check and root redirect if auth failed
func AuthWebsite(action HandlerFunc) HandlerFunc {
	return request(authCheck, failedAuthWebsite, action)
}

// AdminWebsite with admin check and root redirect if auth failed
func AdminWebsite(action HandlerFunc) HandlerFunc {
	return request(adminCheck, failedAuthWebsite, action)
}

// File without auth check, no locale check, and no auth failed function
func File(action HandlerFunc) HandlerFunc {
	return request(noCheck, func(*Connection) {}, action)
}

// WebSocket without auth check
func WebSocket(action WebSocketFunc) HandlerFunc {
	return webSocket(noCheck, func(*Connection) {}, action)
}

// AuthWebSocket with auth check and 401 if auth failed
func AuthWebSocket(action WebSocketFunc) HandlerFunc {
	return webSocket(authCheck, failedAuthRequest, action)
}

// AdminWebSocket with admin check and 401 if auth failed
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
		if con.cfg.InitConnection == nil || con.cfg.Commit == nil || con.cfg.Finish == nil {
			panic(fmt.Sprintf("can't execute request method, some given methods are nil (initConnection: %t, commit: %t, finish: %t)",
				con.cfg.InitConnection == nil,
				con.cfg.Commit == nil,
				con.cfg.Finish == nil))
		}

		if ok := con.cfg.InitConnection(con.server, con); !ok {
			return
		}

		con.Log.Print("Request: ", con.RequestSignature())
		con.Log.Print("Request-Header: ", con.RequestHeader())

		// security check
		if ok := checkAction(con); !ok {
			failAuthAction(con)
		} else {
			action(con)
		}

		response := con.ResponseInfo()
		con.Log.Printf("Response: %s[%d](%s) %s", response.Type, response.Code, response.Duration, response.Description)

		commitChanges := response.Code <= 204 || (response.Code < 400 && response.Type != FailRedirectType)
		con.cfg.Commit(con, commitChanges)

		con.cfg.Finish(con)
	}
}

func webSocket(checkAction CheckFunc, failAuthAction HandlerFunc, action WebSocketFunc) HandlerFunc {
	if checkAction == nil || failAuthAction == nil || action == nil {
		panic(fmt.Sprintf("can't execute request method, some given methods are nil (checkAction: %t, failAuthAction: %t, action: %t)", checkAction == nil, failAuthAction == nil, action == nil))
	}

	return func(con *Connection) {
		if con.cfg.InitConnection == nil || con.cfg.WebSocketError == nil || con.cfg.RefreshWebSocketContext == nil || con.cfg.Commit == nil || con.cfg.Finish == nil {
			panic(fmt.Sprintf("can't execute request method, some given methods are nil (initConnection: %t, webSocketError: %t, refreshWebsocket: %t, commit: %t, finish: %t)",
				con.cfg.InitConnection == nil,
				con.cfg.WebSocketError == nil,
				con.cfg.RefreshWebSocketContext == nil,
				con.cfg.Commit == nil,
				con.cfg.Finish == nil))
		}

		if ok := con.cfg.InitConnection(con.server, con); !ok {
			return
		}

		con.Log.Print("Request: ", con.RequestSignature())
		con.Log.Print("Request-Header: ", con.RequestHeader())

		// security check
		if ok := checkAction(con); !ok {
			failAuthAction(con)

		} else {
			ws, err := con.WebSocketChannel()
			if err != nil {
				con.cfg.WebSocketError(con, err)
				return
			}

			defer ws.Close()

			for {
				msg, open, err := ws.Read()
				if err != nil {
					con.cfg.WebSocketError(con, err)
					return
				}
				if !open {
					break
				}

				err = con.cfg.RefreshWebSocketContext(con.server, con)
				if err != nil {
					con.cfg.WebSocketError(con, err)
					return
				}

				ok := action(msg, con.Context, ws)
				con.cfg.Commit(con, ok)
			}
		}

		response := con.ResponseInfo()
		con.Log.Printf("ResponseInfo: %s[%d](%s) %s", response.Type, response.Code, response.Duration, response.Description)

		con.cfg.Finish(con)
	}
}
