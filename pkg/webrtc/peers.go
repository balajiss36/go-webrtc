package webrtc

import (
	"sync"

	"github.com/balajiss36/go-webrtc/pkg/chat"
	"github.com/fasthttp/websocket"
	"github.com/pion/webrtc/v3"
)

type Room struct {
	Peers Peers
	Hub   *chat.Hub
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Peers struct {
	ListLock    sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	websocket      *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	return t.Conn.WriteJSON(v)
}

func (p *Peers) AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {}

func (p *Peers) RemoveTrack(track *webrtc.TrackLocalStaticRTP) {}

func (p *Peers) SignalPeerConnection() {
}

func (p *Peers) DispatchKeyFrame() {
}
