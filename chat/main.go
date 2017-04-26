package main

import (
	"flag"
	"html/template"
	"log"
	"mycode/trace"
	"net/http"
	"os"
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
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈します
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandlar{filename: "chat.html"}))
	http.Handle("/login", &templateHandlar{filename: "login.html"})
	http.Handle("/room", r)
	go r.run()
	// Webサーバの開始
	log.Println("Webサーバを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
