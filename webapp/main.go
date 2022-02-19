package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// once do は一度だけ呼び出すもの
	t.once.Do(func() {
		// htmlファイルをパース(解析)し、
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// 出力先に渡す。
	t.templ.Execute(w, r)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}

func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	flag.Parse() // parse the flags

	// ルームを作成
	r := newRoom()

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
