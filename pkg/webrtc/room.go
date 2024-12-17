package webrtc

import (
	"log"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/pion/webrtc/v3"
)

func RoomConn(c *websocket.Conn, p *Peers) {
	var config webrtc.Configuration

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {

		log.Fatalf("Error in creating peer connection: %v", err)
		return
	}

	newPeer := PeerConnectionState{
		PeerConnection: peerConnection,
		WebSocket:      &ThreadSafeWriter{},
		Conn:           c,
		Mutex:          sync.Mutex{},
	}
}
