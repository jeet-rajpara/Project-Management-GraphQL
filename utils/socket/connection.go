package socket

import (
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var socketServer *socketio.Server

func SocketConnection() *socketio.Server {

	if socketServer != nil {
		return socketServer
	}

	socketServer = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				Client: &http.Client{
					Timeout: time.Minute,
				},
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})
	return socketServer
}

func GetServer() *socketio.Server {
	return socketServer
}

func CloseServer() {
	if socketServer != nil {
		socketServer.Close()
	}
}
