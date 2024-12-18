package webrtc

import (
	"log"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/pion/webrtc/v3"
)

var (
	RoomsLock sync.RWMutex
	Rooms     map[string]*Room
)

func RoomConn(c *websocket.Conn, p *Peers) {
	var config webrtc.Configuration
	offer := webrtc.SessionDescription{}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {

		log.Fatalf("Error in creating peer connection: %v", err)
		return
	}
	peerConnection.SetRemoteDescription(offer)

	newPeer := PeerConnectionState{
		PeerConnection: peerConnection,
		websocket:      &ThreadSafeWriter{},
		Conn:           c,
		Mutex:          sync.Mutex{},
	}
}
