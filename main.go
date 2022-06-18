package main

import (
	"html/template"
	"net/http"

	"./model"
)

func process(w http.ResponseWriter, r *http.Request) {
	// htmlを解析
	// t, _ := template.ParseFiles("client/index.html")
	// Must関数を挟んでエラーをラップ
	t := template.Must(template.ParseGlob("client/*.html"))
	daysOfWeek := model.Process()
	// daysOfWeek := []string{"月", "火", "水", "木", "金", "土", "日"}
	// dataをテンプレートに当てはめている。
	t.Execute(w, daysOfWeek)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/process", process)

	server.ListenAndServe()
}
