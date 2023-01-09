package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	tr "hub/service"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func receive(c *websocket.Conn, done chan struct{}) {
	defer close(done)

	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("read: ", err)
			break
		}
		fmt.Println("recv: ", string(message))
	}
}

func main() {
	flag.Parse()
	translator := tr.NewTransletor("en", "kk-kz")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/2"}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatalln("dial: ", err)
	}

	defer c.Close()

	done := make(chan struct{})

	//go receive(c, done)
	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				log.Println("Err: ", err)
				break
			}

			res, err := translator.TransletorMSG(string(message))
			if err != nil {
				break
			}
			fmt.Println("Received message: ", res)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print("Message: ")
		text := scanner.Text()
		fmt.Println()
		if text == "" {
			continue
		}
		translator2 := tr.NewTransletor("kk-kz", "en")
		res, _ := translator2.TransletorMSG(text)
		err := c.WriteMessage(websocket.TextMessage, []byte(res))

		if err != nil {
			log.Println("write: ", err)
			break
		}
	}
}
