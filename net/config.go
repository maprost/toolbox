package net

type Config struct {
	// cookie config
	CookieRootPath string
	DefaultCookies []string

	// templates
	StartDelimiter string
	EndDelimiter   string

	// router config
	FailRedirectPath string
	InitConnection   func(*Connection) error
	AuthCheck        func(*Connection) bool
	AdminCheck       func(*Connection) bool
	Close            func(con *Connection, commit bool) error
	Finish           func(con *Connection)
}

func NewConfig() *Config {
	return &Config{
		CookieRootPath:   "/",
		DefaultCookies:   []string{},
		StartDelimiter:   "§§", // works better with html
		EndDelimiter:     "§§", // works better with html
		FailRedirectPath: "/",
		AuthCheck:        nil,
		AdminCheck:       nil,
		InitConnection:   nil,
	}
}
