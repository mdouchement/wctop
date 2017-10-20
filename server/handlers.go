package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mdouchement/wctop/async"
	"github.com/pkg/errors"
)

func homeHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := FSString(localAssets, "/assets/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mw := BuildtimeFilter("text/html", w)
	defer mw.Close()

	async.Start()

	fmt.Fprintf(mw, tmpl)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go echo(conn)
}

func echo(conn *websocket.Conn) {
	fmt.Printf("Remote %s has subscribed\n", conn.RemoteAddr())

	id, notifCh := async.WsNotifier.Subscribe()
	for {
		select {
		case notif := <-notifCh:
			err := conn.WriteJSON(notif)
			if err != nil {
				fmt.Println(errors.Wrap(err, "ws serialization"))
				async.WsNotifier.UnSubscribe(id)
				fmt.Printf("Unsubscribe remote %s\n", conn.RemoteAddr())
				async.Stop() // Stop if no longer subscribers

				err = conn.Close()
				if err = errors.Wrap(err, "ws closing"); err != nil {
					fmt.Println(err)
				}
				return
			}
		}
	}
}
