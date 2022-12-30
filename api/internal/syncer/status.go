package syncer

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

var sockets = make(map[string]*Connection)

type Connection struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

type StatusMessage struct {
	Syncing bool `json:"syncing"`
}

func (c *Connection) Send(message StatusMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Socket.WriteJSON(message)
}

func SendToAll(message StatusMessage) {
	logger := getLogger()
	messageJson, _ := json.Marshal(message)
	logger.Debug("sending message to all clients: ", string(messageJson))
	for key, connection := range sockets {
		logger.Debug("sending to client: ", key)
		err := connection.Send(message)
		if err != nil {
			delete(sockets, key)
		}
	}
}
