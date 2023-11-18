package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

var connections []*websocket.Conn

func remove[T comparable](arr []T, elem T) []T {
	for i, e := range arr {
		if e == elem {
			return slices.Delete[[]T](arr, i, i)
		}
	}
	return arr
}

func closeConn(conn *websocket.Conn) {
	remove(connections, conn)
	conn.Close()
}

func sendmessage(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		closeConn(conn)
		return
	}
	if !slices.Contains[*websocket.Conn](connections, conn) {
		connections = append(connections, conn)
	}
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			closeConn(conn)
			return
		}
		for _, c := range connections {
			err = c.WriteMessage(msgType, msg)
			if err != nil {
				closeConn(conn)
				return
			}
		}

	}

}

func main() {
	srv := http.NewServeMux()
	srv.HandleFunc("/", sendmessage)
	log.Println("Server is starting on 5000 port")
	log.Fatal(http.ListenAndServe(":5000", srv))

}
