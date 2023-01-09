package room

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}

type Connection struct {
	ws   *websocket.Conn
	Send chan []byte
}

func (s Subscription) readPump() {
	c := s.conn

	defer func() {
		H.unregister <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		m := Message{data: msg, room: s.room}
		H.broadcast <- m
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().UTC().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *Subscription) writePump() {
	c := s.conn

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, message)
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func ServerWS(w http.ResponseWriter, r *http.Request, roomId string) {
	fmt.Println("Room ID: ", roomId)

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &Connection{Send: make(chan []byte, 256), ws: ws}

	s := Subscription{c, roomId}

	H.register <- s

	go s.readPump()
	go s.writePump()
}
