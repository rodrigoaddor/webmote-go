package server

import (
	"github.com/pion/webrtc/v3"
	"log"
	"os"
	"strings"
)

type WRTCPeer struct {
	peer *webrtc.PeerConnection
}

func NewWRTCPeer() (*WRTCPeer, error) {
	iceServersURLs := strings.Split(os.Getenv("ICE_SERVERS"), ",")
	iceServers := make([]webrtc.ICEServer, len(iceServersURLs))
	for i, url := range iceServersURLs {
		iceServers[i] = webrtc.ICEServer{
			URLs: []string{url},
		}
	}

	peer, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: iceServers,
	})

	peer.OnDataChannel(func(channel *webrtc.DataChannel) {
		log.Println("got a datachannel")
	})

	return &WRTCPeer{
		peer: peer,
	}, err
}

func (p WRTCPeer) Reply(offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	if err := p.peer.SetRemoteDescription(offer); err != nil {
		return nil, err
	}

	answer, err := p.peer.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	if err := p.peer.SetLocalDescription(answer); err != nil {
		return nil, err
	}
	return &answer, nil
}

func (p WRTCPeer) AddICECandidate(ice webrtc.ICECandidateInit) error {
	return p.peer.AddICECandidate(ice)
}

func (p WRTCPeer) OnICECandidate(f func(candidate *webrtc.ICECandidate)) {
	p.peer.OnICECandidate(f)
}

func (p WRTCPeer) OnDataChannel(f func(channel *webrtc.DataChannel)) {
	p.peer.OnDataChannel(f)
}
