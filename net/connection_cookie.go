package net


func (con *Connection) SaveCookie(name string, value string) {
	con.gin.SetCookie(name, value, 0, con.cfg.CookieRootPath, "", false, false)
}

func (con *Connection) GetCookie(key string) string {
	value, err := con.gin.Cookie(key)
	if err != nil {
		return ""
	}
	return value
}

func (con *Connection) ClearCookie(name string) {
	con.gin.SetCookie(name, "", -1, con.cfg.CookieRootPath, "", false, false)
}

// ================== Default Cookies =====================

func (con *Connection) loadDefaultCookies() {
	for _, key := range con.cfg.DefaultCookies {
		value := con.GetCookie(key)
		con.defaultCookieValues[key] = value
	}
}

func (con *Connection) saveDefaultCookies() {
	for key, value := range con.defaultCookieValues {
		con.SaveCookie(key, value)
	}
}

func (con *Connection) SaveDefaultCookie(name string, value string) {
	con.defaultCookieValues[name] = value
}

func (con *Connection) GetDefaultCookie(key string) string {
	value, ok := con.defaultCookieValues[key]
	if !ok{
		return ""
	}
	return value
}

func (con *Connection) ClearDefaultCookies() {
	for _, key := range con.cfg.DefaultCookies {
		con.ClearCookie(key)
		con.defaultCookieValues[key] = ""
	}
}