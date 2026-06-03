package realtime

import (
	"log"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

const (
	EventJoinAuditStream   = "join_audit_stream"
	EventAuditStreamJoined = "audit_stream_joined"
	EventAuditLogCreated   = "audit_log_created"
	EventAuditLiveError    = "audit_live_error"
)

func (s *Server) registerAuditEvents(client *socket.Socket) {
	client.On(EventJoinAuditStream, func(data ...any) {
		s.handleJoinAuditStream(client, data...)
	})
}

func (s *Server) handleJoinAuditStream(client *socket.Socket, data ...any) {
	teamID, ok := extractTeamID(data...)
	if !ok {
		client.Emit(EventAuditLiveError, map[string]any{
			"message": "teamId is required",
		})
		return
	}

	room := getAuditRoom(teamID)

	client.Join(room)

	client.Emit(EventAuditStreamJoined, map[string]any{
		"teamId": teamID,
		"room":   string(room),
	})

	log.Println("Client joined audit stream:", client.Id(), "team:", teamID)
}

func (s *Server) EmitAuditLogCreated(teamID string, payload any) {
	room := getAuditRoom(teamID)

	s.io.To(room).Emit(EventAuditLogCreated, payload)

	log.Println("Audit log emitted to room:", string(room))
}

func getAuditRoom(teamID string) socket.Room {
	return socket.Room("audit_logs:" + teamID)
}

func extractTeamID(data ...any) (string, bool) {
	if len(data) == 0 {
		return "", false
	}

	payload, ok := data[0].(map[string]any)
	if !ok {
		return "", false
	}

	teamID, ok := payload["teamId"].(string)
	if !ok || teamID == "" {
		return "", false
	}

	return teamID, true
}