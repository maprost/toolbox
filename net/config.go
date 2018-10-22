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
	FailRedirectPath string
	InitContext      func(*Server) (error, interface{})
	AuthCheck        CheckFunc
	AdminCheck       CheckFunc
	WebSocketError   func(con *Connection)
	Close            func(con *Connection, commit bool) error
	Finish           HandlerFunc
}

func NewConfig() *Config {
	return &Config{
		Host:             "localhost",
		Port:             "8080",
		CookieRootPath:   "/",
		DefaultCookies:   []string{},
		StartDelimiter:   "§§", // works better with html
		EndDelimiter:     "§§", // works better with html
		FailRedirectPath: "/",
		InitContext:      nil,
		AuthCheck:        nil,
		AdminCheck:       nil,
		Close:            nil,
		Finish:           nil,
	}
}
