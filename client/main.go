package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var username string = "anon"

func main() {
	u := url.URL{Scheme: "ws", Host: "192.168.3.33:5000", Path: "/echo"}
	log.Printf("\nconnecting to %s", u.String())
	fmt.Println("Use '/setname ...' to set nickname")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("%s\n", message)
		}
	}()
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("> ")
		msg, err := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(msg, "/setname") {
			username = strings.Replace(msg, "/setname ", "", -1)
			username = strings.TrimSpace(username)
		} else {
			fmt.Print("\033[F")
			msg = username + " : " + msg
			c.WriteMessage(2, []byte(msg))
		}
	}

}
