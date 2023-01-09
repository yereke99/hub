package room

type Message struct {
	data []byte
	room string
}

type Subscription struct {
	conn *Connection
	room string
}

type Hub struct {
	rooms map[string]map[*Connection]bool

	broadcast chan Message

	register chan Subscription

	unregister chan Subscription
}

var H = Hub{
	rooms:      make(map[string]map[*Connection]bool),
	broadcast:  make(chan Message),
	register:   make(chan Subscription),
	unregister: make(chan Subscription),
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.room]

			if _, ok := connections[s.conn]; ok {
				delete(connections, s.conn)
				close(s.conn.Send)
				if len(connections) == 0 {
					delete(h.rooms, s.room)
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.Send <- m.data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
