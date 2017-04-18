package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// temp1は1つのテンプレート
type templateHandlar struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理する
func (t *templateHandlar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandlar{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	// Webサーバの開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
