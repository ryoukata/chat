package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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
	r := newRoom()
	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// start ChatRoom
	go r.run()
	// start Web Server
	log.Println("Start WebServer. Port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
