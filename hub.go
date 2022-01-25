package main

type Hub struct {
	// クライアントの登録
	clients map[*Client]bool

	// クライアントからのメッセージを受信
	broadcast chan []byte

	// クライアントからのリクエストを登録
	register chan *Client

	// クライアントからのリクエスト解除
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		// 入室
		case client := <-h.register:
			h.clients[client] = true
		// 退室
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// 全員に送信
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				// 送信
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
