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

func (s *Server) registerInternalDirectMessageRoutes() {
	s.mux.HandleFunc("/internal/direct-messages/emit-created", s.handleEmitDirectMessageCreated)
	s.mux.HandleFunc("/internal/direct-messages/emit-reaction-updated", s.handleEmitDirectMessageReactionUpdated)
	s.mux.HandleFunc("/internal/direct-messages/emit-deleted", s.handleEmitDirectMessageDeleted)
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "direct message deleted emitted",
	})
}