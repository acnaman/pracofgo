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

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
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

	//_apikey := "AIzaSyBAqJ5bYibVfvkIyo8xkgocWSnKhdBffmw"
	_secureKey := "hogehogefugafuga"
	_clientid := "21664827745-jgjl7idpso07kkaek9eb0rbijtuvlos2.apps.googleusercontent.com"
	_secret := "3Yzyss7hR5poI5qBeUkFQCHC"
	_callbackFacebook := "http://localhost:8080/auth/callback/facebook"
	_callbackGithub := "http://localhost:8080/auth/callback/github"
	_callbackGoogle := "http://localhost:8080/auth/callback/google"

	flag.Parse() // フラグを解釈します
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey(_secureKey)
	gomniauth.WithProviders(
		facebook.New(_clientid, _secret, _callbackFacebook),
		github.New(_clientid, _secret, _callbackGithub),
		google.New(_clientid, _secret, _callbackGoogle),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandlar{filename: "chat.html"}))
	http.Handle("/login", &templateHandlar{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	// Webサーバの開始
	log.Println("Webサーバを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
