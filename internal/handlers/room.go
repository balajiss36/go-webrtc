package handlers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/balajiss36/go-webrtc/pkg/chat"
	w "github.com/balajiss36/go-webrtc/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func CreateRoom(c *fiber.Ctx) error {
	roomUUID := uuid.New().String()
	fmt.Println(roomUUID)
	return c.Status(http.StatusOK).JSON(fiber.Map{"roomID": roomUUID})
}

func GetRoom(c *fiber.Ctx) error {
	uuid := c.Path("roomID")
	if uuid == "" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "roomID is required"})
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "production" {
		ws = "wss"
	}

	uuid, suuid, _ := createOrGetRoom(uuid)

	fmt.Println(uuid, suuid)

	fmt.Println(ws)

	roomDetails := map[string]string{
		"RoomWebSocketAddr":   fmt.Sprintf("%s://%s/room/%s/websocket", ws, c.Hostname(), uuid),
		"RoomLink":            fmt.Sprintf("%s://%s/room/%s", c.Context().URI().Scheme(), c.Hostname(), uuid),
		"ChatWebSocketAddr":   fmt.Sprintf("%s://%s/room/%s/chat/websocket", ws, c.Hostname(), uuid),
		"ViewerWebSocketAddr": fmt.Sprintf("%s://%s/room/%s/viewer/websocket", ws, c.Hostname(), uuid),
		"StreamLink":          fmt.Sprintf("%s://%s/stream/%s", c.Context().URI().Scheme(), c.Hostname(), suuid),
		"Type":                "room",
	}
	return c.Status(http.StatusOK).JSON(roomDetails)
}

func RoomWebsocket(c *websocket.Conn) {
	uuid := c.Params("roomID")
	if uuid == "" {
		log.Printf("roomID is required")
		return
	}

	_, _, room := createOrGetRoom(uuid)
	w.RoomConn(c, room.Peers)
}

func createOrGetRoom(uuid string) (string, string, *w.Room) {
	w.RoomsLock.Lock()
	defer w.RoomsLock.Unlock()

	h := sha256.New()
	h.Write([]byte(uuid))
	// secure uuid
	suuid := fmt.Sprintf("%x", h.Sum(nil))
	if room := w.Rooms[uuid]; room != nil {
		if _, ok := w.Streams[suuid]; !ok {
			w.Streams[suuid] = room
		}
		return uuid, suuid, room
	}

	hub := chat.NewHub()
	p := &w.Peers{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &w.Room{
		Peers: p,
		Hub:   hub,
	}
	w.Rooms[uuid] = room
	w.Streams[suuid] = room
	go hub.Run()
	fmt.Println(room)
	return uuid, "", room
}

func RoomViewerWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}
	w.RoomsLock.Lock()
	defer w.RoomsLock.Unlock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomsLock.Unlock()
		roomViewerConn(c, peer.Peers)
		return
	}
}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	defer c.Close()
	for {
		select {
		case <-ticker.C:
			w, err := c.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}
