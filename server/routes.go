package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pdk/kickstart/assets"
)

func (s Server) makeRouter() *mux.Router {

	r := mux.NewRouter()

	// serve static assets
	r.PathPrefix("/css/").Handler(http.FileServer(http.FS(assets.FS)))
	r.PathPrefix("/img/").Handler(http.FileServer(http.FS(assets.FS)))
	r.PathPrefix("/js/").Handler(http.FileServer(http.FS(assets.FS)))

	r.Path("/").HandlerFunc(writeIndex)

	r.PathPrefix("/").HandlerFunc(write404)

	return r
}

var (
	indexHTML, notFoundHTML string
)

func init() {
	var err error

	indexHTML, err = b2s(assets.FS.ReadFile("templates/index.html"))
	if err != nil {
		log.Fatalf("cannot read templates/index.html")
	}

	notFoundHTML, err = b2s(assets.FS.ReadFile("templates/404.html"))
	if err != nil {
		log.Fatalf("cannot read templates/404.html")
	}
}

func writeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, indexHTML)
}

func write404(w http.ResponseWriter, r *http.Request) {

	log.Printf("NotFound: %s", r.RequestURI)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprintln(w, notFoundHTML)
}

// b2s quick convert []byte to string (with error passthru)
func b2s(b []byte, err error) (string, error) {
	return string(b), err
}
