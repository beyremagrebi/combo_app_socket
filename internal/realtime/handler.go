package realtime

import (
	"log"

	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (s *Server) registerHandlers() {
	s.io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		s.handleConnection(client)
	})
}

func (s *Server) handleConnection(client *socket.Socket) {
	log.Println("Client connected:", client.Id())

	client.Emit("server_message", "Welcome from Go Socket.IO server")

	s.registerAuditEvents(client)
	s.registerMessageEvents(client)

	client.On("disconnect", func(reason ...any) {
		s.handleDisconnect(client, reason...)
	})
}

func (s *Server) handleDisconnect(client *socket.Socket, reason ...any) {
	log.Println("Client disconnected:", client.Id(), reason)
}