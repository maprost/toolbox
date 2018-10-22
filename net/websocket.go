package net

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// WebSocketChannel
// example:
//	ws, err := con.WebSocketChannel()
//	if err != nil {
//		return err
//	}
//
//	defer ws.Close()
//
//	for {
//		msg, err := ws.Read()
//		if err != nil {
//			return err
// 		}
// 		// implement message handling
//	}
//}
type WebSocketChannel struct {
	conn *websocket.Conn
}

func (ws *WebSocketChannel) Read() ([]byte, bool, error) {
	msgType, msg, err := ws.conn.ReadMessage()
	fmt.Println("MsgType:", msgType, string(msg))
	if _, ok := err.(*websocket.CloseError); ok {
		return msg, false, nil
	}

	return msg, msgType > 0, err
}

func (ws *WebSocketChannel) ReadJson(js interface{}) (bool, error) {
	err := ws.conn.ReadJSON(js)
	if _, ok := err.(*websocket.CloseError); ok {
		return false, nil
	}
	return true, err
}

func (ws *WebSocketChannel) WriteText(data string) error {
	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(data))
	return err
}

func (ws *WebSocketChannel) WriteBinary(data []byte) error {
	err := ws.conn.WriteMessage(websocket.BinaryMessage, data)
	return err
}

func (ws *WebSocketChannel) WriteJson(js interface{}) error {
	err := ws.conn.WriteJSON(js)
	return err
}

func (ws *WebSocketChannel) Close() error {
	err := ws.conn.WriteMessage(websocket.CloseMessage, []byte{3, 232})
	if err != nil {
		return err
	}

	return ws.conn.Close()
}
