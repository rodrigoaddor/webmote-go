package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
	"github.com/rodrigoaddor/webmote-go/pkg/data/gen"
	"github.com/rodrigoaddor/webmote-go/pkg/server"
	"github.com/rodrigoaddor/webmote-go/pkg/utils"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ws", func(ctx *gin.Context) {
		if err := server.HandleWS[*gen.SignalingMessage](ctx, handleWebRTC); err != nil {
			_ = ctx.AbortWithError(500, err)
		}
	})

	return r
}

func handleWebRTC(ctx *gin.Context, data *gen.SignalingMessage) {
	peer, hasPeer := utils.Get[*server.WRTCPeer](ctx, "peer")
	conn, hasConn := utils.Get[*websocket.Conn](ctx, "ws")
	if !hasConn {
		_ = ctx.AbortWithError(500, errors.New("no websocket connection in context"))
		return
	}

	if !hasPeer {
		newPeer, err := server.NewWRTCPeer()

		//newPeer.OnDataChannel(func(channel *webrtc.DataChannel) {
		//	channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		//		log.Printf("> %v", msg.Data)
		//	})
		//})

		newPeer.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			data := candidate.ToJSON()

			err := wspb.Write(ctx, *conn, &gen.SignalingMessage{
				Signaling: &gen.SignalingMessage_IceCandidate{
					IceCandidate: &gen.IceCandidate{
						Candidate:        data.Candidate,
						SdpMid:           *data.SDPMid,
						SdpMLineIndex:    uint32(*data.SDPMLineIndex),
						UsernameFragment: *data.UsernameFragment,
					},
				},
			})

			if err != nil {
				_ = ctx.AbortWithError(500, err)
				return
			}
		})

		peer = &newPeer
		if err != nil {
			_ = ctx.AbortWithError(500, err)
			return
		}

		ctx.Set("peer", newPeer)
	}

	if negotiation := data.GetNegotiation(); negotiation != nil {
		answer, err := (*peer).Reply(webrtc.SessionDescription{
			Type: webrtc.SDPTypeOffer,
			SDP:  negotiation.Sdp,
		})
		if err != nil {
			_ = ctx.AbortWithError(500, err)
			return
		}

		err = wspb.Write(ctx, *conn, &gen.SignalingMessage{
			Signaling: &gen.SignalingMessage_Negotiation{
				Negotiation: &gen.Negotiation{
					Type: "answer",
					Sdp:  answer.SDP,
				},
			},
		})

		if err != nil {
			_ = ctx.AbortWithError(500, err)
			return
		}

	}

	if ice := data.GetIceCandidate(); ice != nil {

	}
}
