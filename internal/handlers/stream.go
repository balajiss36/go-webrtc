package handlers

import (
	"fmt"
	"os"
	"time"

	w "github.com/balajiss36/go-webrtc/pkg/webrtc"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
)

func Stream(c *fiber.Ctx) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		return c.Status(fiber.StatusBadRequest).SendString("suuid is required")
	}
	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "production" {
		ws = "wss"
	}
	w.RoomsLock.Lock()
	if _, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebSocketAddr": fmt.Sprintf("%s://%s/stream/%s/websocket", ws, c.Hostname(), suuid),
			"StreamLink":          fmt.Sprintf("%s://%s/stream/%s", c.Protocol(), c.Hostname(), suuid),
			"ChatWebSocketAddr":   fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Hostname(), suuid),
			"ViewerWebSockerAddr": fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Hostname(), suuid),
			"Type":                "stream",
		}, "layouts/main")
	}
	w.RoomsLock.Unlock()

	return c.Render("stream", fiber.Map{"NoStream": "true", "Leave": true}, "layouts/main")
}

func StreamWebsocket(c *websocket.Conn) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		return nil
	}
	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		w.StreamConn(c, stream.Peers)
		return
	}
	w.RoomsLock.Unlock()
}

func StreamViewerWebSocket(c *websocket.Conn) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		return nil
	}
	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		viewerConn(c, stream.Peers)
		return
	}
	w.RoomsLock.Unlock()
}

func viewerConn(c *websocket.Conn, p *w.Peers) error {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	defer c.Close()
	for {
		select {
		case <-ticker.C:
			w, err := c.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}
