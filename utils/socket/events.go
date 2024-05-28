package socket

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func RegisterEvents(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		s.Join("my-room")
		return nil
	})

	// server.OnEvent("/", "share_project", func(s socketio.Conn, msg string) {
	// 	log.Println("share_project event:", msg)
	// 	s.Emit("reply", "received "+msg)
	// })

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("disconnected:", s.ID(), reason)
	})
}
