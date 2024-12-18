package handlers

import (
	"github.com/balajiss36/go-webrtc/pkg/chat"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	w "github.com/balajiss36/go-webrtc/pkg/webrtc"
)

func RoomChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{}, "layouts/main")
}

func RoomChatWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}
	w.RoomsLock.Lock()

	room := w.Rooms[uuid]
	if room == nil {
		return
	}

	chat.PeerChatConn(c.Conn, room.Hub)
}

func StreamChatWebsocket(c *websocket.Conn) {
	suuid := c.Params("suuid")
	if suuid == "" {
		return
	}
	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		chat.PeerChatConn(c.Conn, stream.Hub)
		return
	}
}
