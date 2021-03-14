package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/ryoukata/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

// templ return one template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServerHttp execute HTTP Request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Application Address")
	flag.Parse()
	// set up Gomniauth
	gomniauth.SetSecurityKey("2hi0oeq67ao05kg9a")
	gomniauth.WithProviders(
		facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
		github.New("", "", "http://localhost:8080/auth/callback/github"),
		google.New("439060568234-hggfg4uqkf9gkhp8djdrq140p1607g0k.apps.googleusercontent.com", "fMjlDA04FkN-vswif2zRTBFM", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// root
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// start ChatRoom
	go r.run()
	// start Web Server
	log.Println("Start WebServer. Port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
