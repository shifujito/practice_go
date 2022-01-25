package main

import "github.com/gorilla/websocket"

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	// クライアントがwebsocketからreadmessageを使ってデータを読み込む
	for {
		// データの読み込みを行う
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			// roomのforwardチャネルに送る
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	// sendチャネルからメッセージを受け取り、websocketのwritemessageを使って書き出す。
	for msg := range c.send {
		// websocketへの書き込み
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
