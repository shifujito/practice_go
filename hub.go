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
