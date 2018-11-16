package net

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// PathParam returns the value of the given path key, "" if the key doesn't exists
func (con *Connection) PathParam(key string) string {
	return con.gin.Param(key)
}

// QueryParam returns the value of the given query key, "" if the key doesn't exists
func (con *Connection) QueryParam(key string) string {
	return con.gin.Query(key)
}

// QueryParamArray returns the values of the given query key, empty list if the key doesn't exists
func (con *Connection) QueryParamArray(key string) []string {
	return con.gin.QueryArray(key)
}

// DefaultQueryParam returns the value of the given query key, 'defaultValue' if the key doesn't exists
func (con *Connection) DefaultQueryParam(key string, defaultValue string) string {
	return con.gin.DefaultQuery(key, defaultValue)
}

// QueryParamDate returns the time value of the given query key, 'time.Time{}' if the key doesn't exists
func (con *Connection) QueryParamDate(key string, dateFormat string) (time.Time, error) {
	q := con.gin.Query(key)

	date, err := time.Parse(dateFormat, q)
	if err != nil {
		return time.Time{}, NewNotFoundError("can't convert query param '%s'-'%v' in date format '%s'.", key, q, dateFormat)
	}

	return date, nil
}

// Body converts the body into the given 'obj'
func (con *Connection) Body(obj interface{}) error {
	bytes, err := ioutil.ReadAll(con.gin.Request.Body)
	if err != nil {
		return NewNotFoundError("can't read body")
	}

	fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return NewNotFoundError("can't convert body (%s)", string(bytes))
	}
	return nil
}

func (con *Connection) WebSocket(messageChannel chan []byte) error {
	ws, err := con.WebSocketChannel()
	if err != nil {
		return err
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	for {
		// Read in a new message
		msg, open, err := ws.Read()
		if err != nil {
			return err
		}
		if !open {
			break
		}
		messageChannel <- msg
	}

	return nil
}

func (con *Connection) WebSocketChannel() (*WebSocketChannel, error) {
	// Upgrade initial GET request to a websocket
	ws, err := con.wsUpgrader.Upgrade(con.gin.Writer, con.gin.Request, nil)
	return &WebSocketChannel{conn: ws}, err
}
