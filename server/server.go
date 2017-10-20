package server

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/bmizerany/pat"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

var localAssets bool

func init() {
	localAssets = os.Getenv("LOCAL_ASSETS") != ""
	if localAssets {
		fmt.Println("Local assets activated")
	} else {
		fmt.Println("BIN-DATA assets activated")
	}
}

// Run starts the server on the given socket.
func Run(socket string) error {
	m := pat.New()
	m.Get("/", http.HandlerFunc(homeHandler))

	http.Handle("/", m)
	http.Handle("/assets/", asset())
	http.HandleFunc("/ws", wsHandler)
	fmt.Printf("== Starting server on http://%s\n", socket)
	return http.ListenAndServe(socket, nil)
}

func asset() http.Handler {
	fs := http.FileServer(FS(localAssets))
	if localAssets {
		return fs
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return m.Middleware(fs)
}
