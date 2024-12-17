package chat

import (
	"bytes"
	"log"
	"time"

	"github.com/fasthttp/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingWait       = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Fatalf("Read deadline exceeded: %v", err)
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Fatalf("set read deadline %v", err)
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Hub.broadcast <- message

	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(5 * time.Minute)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Fatalf("set write deadline %v", err)
				return
			}
			if !ok {
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatalf("new writer error while sending textmessage %v", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Fatalf("writer error %v", err)
				return
			}
			n := len(c.Send)
			// for multiple lines of messages, iterate over the channel after newline and write to the conn
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					log.Fatalf("writer error while sending newline %v", err)
					return
				}
				_, err = w.Write(<-c.Send)
				if err != nil {
					log.Fatalf("writer error while sending next line %v", err)
					return
				}
			}
			if err := w.Close(); err != nil {
				return
			}
			// periodically send ping messages to channel to keep the conn alive
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Fatalf("writer error while setting deadline for ticket %v", err)

				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Fatalf("writer error while sending ping message %v", err)
				return
			}
		}
	}
}

func PeerChatConn(c *websocket.Conn, h *Hub) {
	client := &Client{
		Hub:  h,
		Conn: c,
		Send: make(chan []byte, 256),
	}
	client.Hub.register <- client

	go client.WritePump()
	client.ReadPump()
}
