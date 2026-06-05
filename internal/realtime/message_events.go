package realtime

import (
	"log"

	"github.com/zishang520/socket.io/servers/socket/v3"
)

const (
	EventJoinMessageStream      = "join_message_stream"
	EventLeaveMessageStream     = "leave_message_stream"
	EventMessageStreamJoined    = "message_stream_joined"
	EventMessageCreated         = "message_created"
	EventMessageReactionUpdated = "message_reaction_updated"
	EventMessageLiveError       = "message_live_error"
	EventMessageDeleted         = "message_deleted"

	EventJoinDirectMessageStream      = "join_direct_message_stream"
	EventLeaveDirectMessageStream     = "leave_direct_message_stream"
	EventDirectMessageStreamJoined    = "direct_message_stream_joined"
	EventDirectMessageCreated         = "direct_message_created"
	EventDirectMessageReactionUpdated = "direct_message_reaction_updated"
	EventDirectMessageDeleted         = "direct_message_deleted"

	EventJoinTeamStream                = "join_team_stream"
	EventLeaveTeamStream               = "leave_team_stream"
	EventTeamStreamJoined              = "team_stream_joined"
	EventUserStatusUpdated             = "user_status_updated"
	EventWorkspaceUnreadMessageCreated = "workspace_unread_message_created"
)

func (s *Server) registerMessageEvents(client *socket.Socket) {
	client.On(EventJoinMessageStream, func(data ...any) {
		s.handleJoinMessageStream(client, data...)
	})

	client.On(EventLeaveMessageStream, func(data ...any) {
		s.handleLeaveMessageStream(client, data...)
	})

	client.On(EventJoinDirectMessageStream, func(data ...any) {
		s.handleJoinDirectMessageStream(client, data...)
	})

	client.On(EventLeaveDirectMessageStream, func(data ...any) {
		s.handleLeaveDirectMessageStream(client, data...)
	})

	client.On(EventJoinTeamStream, func(data ...any) {
		s.handleJoinTeamStream(client, data...)
	})

	client.On(EventLeaveTeamStream, func(data ...any) {
		s.handleLeaveTeamStream(client, data...)
	})

	client.On(EventUserStatusUpdated, func(data ...any) {
		s.handleUserStatusUpdated(client, data...)
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

func (s *Server) handleJoinDirectMessageStream(client *socket.Socket, data ...any) {
	conversationID, ok := extractConversationID(data...)
	if !ok {
		client.Emit(EventMessageLiveError, map[string]any{
			"message": "conversationId is required",
		})
		return
	}

	room := getDirectMessageRoom(conversationID)

	client.Join(room)

	client.Emit(EventDirectMessageStreamJoined, map[string]any{
		"conversationId": conversationID,
		"room":           string(room),
	})

	log.Println(
		"Client joined direct message stream:",
		client.Id(),
		"conversation:",
		conversationID,
	)
}

func (s *Server) handleLeaveDirectMessageStream(client *socket.Socket, data ...any) {
	conversationID, ok := extractConversationID(data...)
	if !ok {
		return
	}

	room := getDirectMessageRoom(conversationID)

	client.Leave(room)

	log.Println(
		"Client left direct message stream:",
		client.Id(),
		"conversation:",
		conversationID,
	)
}

func (s *Server) handleJoinTeamStream(client *socket.Socket, data ...any) {
	teamID, ok := extractTeamID(data...)
	if !ok {
		client.Emit(EventMessageLiveError, map[string]any{
			"message": "teamId is required",
		})
		return
	}

	room := getMessageTeamRoom(teamID)

	client.Join(room)

	client.Emit(EventTeamStreamJoined, map[string]any{
		"teamId": teamID,
		"room":   string(room),
	})

	log.Println("Client joined team stream:", client.Id(), "team:", teamID)
}

func (s *Server) handleLeaveTeamStream(client *socket.Socket, data ...any) {
	teamID, ok := extractTeamID(data...)
	if !ok {
		return
	}

	room := getMessageTeamRoom(teamID)

	client.Leave(room)

	log.Println("Client left team stream:", client.Id(), "team:", teamID)
}

func (s *Server) handleUserStatusUpdated(client *socket.Socket, data ...any) {
	if len(data) == 0 {
		return
	}

	payload, ok := data[0].(map[string]any)
	if !ok {
		return
	}

	teamID, ok := payload["teamId"].(string)
	if !ok || teamID == "" {
		client.Emit(EventMessageLiveError, map[string]any{
			"message": "teamId is required",
		})
		return
	}

	room := getMessageTeamRoom(teamID)

	s.io.To(room).Emit(EventUserStatusUpdated, payload)

	log.Println("User status emitted to team room:", string(room))
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
		"messageId": messageID,
		"channelId": channelID,
	})

	log.Println("Message delete emitted to room:", string(room))
}

func (s *Server) EmitDirectMessageCreated(conversationID string, message map[string]any) {
	room := getDirectMessageRoom(conversationID)

	s.io.To(room).Emit(EventDirectMessageCreated, map[string]any{
		"conversationId": conversationID,
		"message":        message,
	})

	log.Println("Direct message emitted to room:", string(room))
}

func (s *Server) EmitDirectMessageReactionUpdated(conversationID string, message map[string]any) {
	room := getDirectMessageRoom(conversationID)

	s.io.To(room).Emit(EventDirectMessageReactionUpdated, map[string]any{
		"conversationId": conversationID,
		"message":        message,
	})

	log.Println("Direct message reaction emitted to room:", string(room))
}

func (s *Server) EmitDirectMessageDeleted(conversationID string, messageID string) {
	room := getDirectMessageRoom(conversationID)

	s.io.To(room).Emit(EventDirectMessageDeleted, map[string]any{
		"id":             messageID,
		"messageId":      messageID,
		"conversationId": conversationID,
	})

	log.Println("Direct message delete emitted to room:", string(room))
}

func (s *Server) EmitWorkspaceUnreadMessageCreated(teamID string, payload any) {
	room := getMessageTeamRoom(teamID)

	s.io.To(room).Emit(EventWorkspaceUnreadMessageCreated, payload)

	log.Println("Workspace unread message emitted to team room:", string(room))
}

func getMessageRoom(channelID string) socket.Room {
	return socket.Room("channel_messages:" + channelID)
}

func getDirectMessageRoom(conversationID string) socket.Room {
	return socket.Room("direct_messages:" + conversationID)
}

func getMessageTeamRoom(teamID string) socket.Room {
	return socket.Room("team:" + teamID)
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

func extractConversationID(data ...any) (string, bool) {
	if len(data) == 0 {
		return "", false
	}

	payload, ok := data[0].(map[string]any)
	if !ok {
		return "", false
	}

	conversationID, ok := payload["conversationId"].(string)
	if !ok || conversationID == "" {
		return "", false
	}

	return conversationID, true
}