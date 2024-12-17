package handlers

import (
	"github.com/balajiss36/go-webrtc/pkg/chat"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"

	w "github.com/balajiss36/go-webrtc/pkg/webrtc"
)

func RoomChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{}, "layouts/main")
}

func RoomChatWebsocker(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomsLock.Lock()

	room := w.Rooms[uuid]
	if room == nil {
		return
	}

	chat.PeerChatConn(c, room.Hub)
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
		chat.PeerChatConn(c, stream.Hub)
		return
	}
}
