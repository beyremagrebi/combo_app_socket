package realtime

import (
	"log"

	"github.com/zishang520/socket.io/servers/socket/v3"
)

const (
	EventJoinMessageStream     = "join_message_stream"
	EventLeaveMessageStream    = "leave_message_stream"
	EventMessageStreamJoined   = "message_stream_joined"
	EventMessageCreated        = "message_created"
	EventMessageReactionUpdated = "message_reaction_updated"
	EventMessageLiveError      = "message_live_error"
	EventMessageDeleted = "message_deleted"
)

func (s *Server) registerMessageEvents(client *socket.Socket) {
	client.On(EventJoinMessageStream, func(data ...any) {
		s.handleJoinMessageStream(client, data...)
	})

	client.On(EventLeaveMessageStream, func(data ...any) {
		s.handleLeaveMessageStream(client, data...)
	})
}

func (s *Server) handleJoinMessageStream(client *socket.Socket, data ...any) {
	channelID, ok := extractChannelID(data...)
	if !ok {
		client.Emit(EventMessageLiveError, map[string]any{
			"message": "channelId is required",
		})
		return
	}

	room := getMessageRoom(channelID)

	client.Join(room)

	client.Emit(EventMessageStreamJoined, map[string]any{
		"channelId": channelID,
		"room":      string(room),
	})

	log.Println("Client joined message stream:", client.Id(), "channel:", channelID)
}

func (s *Server) handleLeaveMessageStream(client *socket.Socket, data ...any) {
	channelID, ok := extractChannelID(data...)
	if !ok {
		return
	}

	room := getMessageRoom(channelID)

	client.Leave(room)

	log.Println("Client left message stream:", client.Id(), "channel:", channelID)
}

func (s *Server) EmitMessageCreated(channelID string, payload any) {
	room := getMessageRoom(channelID)

	s.io.To(room).Emit(EventMessageCreated, payload)

	log.Println("Message emitted to room:", string(room))
}

func (s *Server) EmitMessageReactionUpdated(channelID string, payload any) {
	room := getMessageRoom(channelID)

	s.io.To(room).Emit(EventMessageReactionUpdated, payload)

	log.Println("Message reaction emitted to room:", string(room))
}

func (s *Server) EmitMessageDeleted(channelID string, messageID string) {
	room := getMessageRoom(channelID)

	s.io.To(room).Emit(EventMessageDeleted, map[string]any{
		"id":        messageID,
		"channelId": channelID,
	})

	log.Println("Message delete emitted to room:", string(room))
}

func getMessageRoom(channelID string) socket.Room {
	return socket.Room("channel_messages:" + channelID)
}

func extractChannelID(data ...any) (string, bool) {
	if len(data) == 0 {
		return "", false
	}

	payload, ok := data[0].(map[string]any)
	if !ok {
		return "", false
	}

	channelID, ok := payload["channelId"].(string)
	if !ok || channelID == "" {
		return "", false
	}

	return channelID, true
}