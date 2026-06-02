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

	client.Emit("server_message", "Welcome from Go Socket.IO v4 server")

	client.On("chat_message", func(data ...any) {
		s.handleChatMessage(client, data...)
	})

	client.On("disconnect", func(reason ...any) {
		s.handleDisconnect(client, reason...)
	})
}

func (s *Server) handleChatMessage(client *socket.Socket, data ...any) {
	log.Println("chat_message:", data)

	client.Emit("chat_message", data...)
	s.io.Emit("new_message", data...)
}

func (s *Server) handleDisconnect(client *socket.Socket, reason ...any) {
	log.Println("Client disconnected:", client.Id(), reason)
}