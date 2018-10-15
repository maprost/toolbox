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
	InitConnection   func(*Server, *Connection) error
	AuthCheck        CheckFunc
	AdminCheck       CheckFunc
	Close            func(con *Connection, commit bool) error
	Finish           func(con *Connection)
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
		InitConnection:   nil,
		AuthCheck:        nil,
		AdminCheck:       nil,
		Close:            nil,
		Finish:           nil,
	}
}
