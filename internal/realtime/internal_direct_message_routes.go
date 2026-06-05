package realtime

import (
	"encoding/json"
	"net/http"
)

type EmitDirectMessageRequest struct {
	ConversationID string         `json:"conversationId"`
	Message        map[string]any `json:"message"`
}

type EmitDirectMessageDeletedRequest struct {
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
}

type EmitWorkspaceUnreadMessageCreatedRequest struct {
	TeamID  string         `json:"teamId"`
	Payload map[string]any `json:"payload"`
}

func (s *Server) registerInternalDirectMessageRoutes() {
	s.mux.HandleFunc(
		"/internal/direct-messages/emit-created",
		s.handleEmitDirectMessageCreated,
	)

	s.mux.HandleFunc(
		"/internal/direct-messages/emit-reaction-updated",
		s.handleEmitDirectMessageReactionUpdated,
	)

	s.mux.HandleFunc(
		"/internal/direct-messages/emit-deleted",
		s.handleEmitDirectMessageDeleted,
	)

	s.mux.HandleFunc(
		"/internal/workspace-unread/emit-message-created",
		s.handleEmitWorkspaceUnreadMessageCreated,
	)
}

func (s *Server) handleEmitDirectMessageCreated(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitDirectMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.ConversationID == "" {
		http.Error(w, "conversationId is required", http.StatusBadRequest)
		return
	}

	if body.Message == nil {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	s.EmitDirectMessageCreated(body.ConversationID, body.Message)

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "direct message emitted",
	})
}

func (s *Server) handleEmitDirectMessageReactionUpdated(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitDirectMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.ConversationID == "" {
		http.Error(w, "conversationId is required", http.StatusBadRequest)
		return
	}

	if body.Message == nil {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	s.EmitDirectMessageReactionUpdated(body.ConversationID, body.Message)

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "direct message reaction emitted",
	})
}

func (s *Server) handleEmitDirectMessageDeleted(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitDirectMessageDeletedRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.ConversationID == "" {
		http.Error(w, "conversationId is required", http.StatusBadRequest)
		return
	}

	if body.MessageID == "" {
		http.Error(w, "messageId is required", http.StatusBadRequest)
		return
	}

	s.EmitDirectMessageDeleted(body.ConversationID, body.MessageID)

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "direct message deleted emitted",
	})
}

func (s *Server) handleEmitWorkspaceUnreadMessageCreated(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitWorkspaceUnreadMessageCreatedRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.TeamID == "" {
		http.Error(w, "teamId is required", http.StatusBadRequest)
		return
	}

	if body.Payload == nil {
		http.Error(w, "payload is required", http.StatusBadRequest)
		return
	}

	s.EmitWorkspaceUnreadMessageCreated(body.TeamID, body.Payload)

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "workspace unread message emitted",
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}