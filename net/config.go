package net

type Config struct {
	// server
	Host string
	Port string

	// cookie config
	CookieRootPath string
	DefaultCookies []string

	// templates
	StartDelimiter string
	EndDelimiter   string

	// handler func
	FailRedirectPath        string
	InitConnection          func(server *Server, con *Connection) (ok bool)   // init context and logger
	RefreshWebSocketContext func(server *Server, con *Connection) (err error) // important for web sockets
	AuthCheck               CheckFunc
	AdminCheck              CheckFunc
	WebSocketError          func(con *Connection, err error)
	Commit                  func(con *Connection, commit bool)
	Finish                  HandlerFunc
}

func NewConfig() *Config {
	return &Config{
		Host:                    "localhost",
		Port:                    "8080",
		CookieRootPath:          "/",
		DefaultCookies:          []string{},
		StartDelimiter:          "§§", // works better with html
		EndDelimiter:            "§§", // works better with html
		FailRedirectPath:        "/",
		InitConnection:          nil,
		AuthCheck:               nil,
		AdminCheck:              nil,
		WebSocketError:          nil,
		RefreshWebSocketContext: nil,
		Commit:                  nil,
		Finish:                  nil,
	}
}
