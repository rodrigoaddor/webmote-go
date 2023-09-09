package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"
	"reflect"
)

type WSHandler[K proto.Message] func(ctx *gin.Context, data K)

func HandleWS[K proto.Message](ctx *gin.Context, handler WSHandler[K]) error {
	conn, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
		Subprotocols:       []string{"webmote"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}

	ctx.Set("ws", conn)

	for {
		var data K
		{
			typ := reflect.TypeOf(data)
			if typ.Kind() != reflect.Pointer {
				return errors.New("K is not a pointer")
			}
			elm := typ.Elem()
			if elm.Kind() != reflect.Struct {
				return errors.New("K is not a struct pointer")
			}
			data = reflect.New(elm).Interface().(K)
		}
		if err := wspb.Read(ctx, conn, data); err != nil {
			return err
		}

		handler(ctx, data)
	}
}
