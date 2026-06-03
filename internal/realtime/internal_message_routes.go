package realtime

import (
	"encoding/json"
	"net/http"
)

type EmitMessageRequest struct {
	ChannelID string         `json:"channelId"`
	Message   map[string]any `json:"message"`
}

type EmitMessageDeletedRequest struct {
	ChannelID string `json:"channelId"`
	MessageID string `json:"messageId"`
}

func (s *Server) registerInternalMessageRoutes() {
	s.mux.HandleFunc("/internal/messages/emit-created", s.handleEmitMessageCreated)
	s.mux.HandleFunc("/internal/messages/emit-reaction-updated", s.handleEmitMessageReactionUpdated)
	s.mux.HandleFunc("/internal/messages/emit-deleted", s.handleEmitMessageDeleted)
}

func (s *Server) handleEmitMessageCreated(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.ChannelID == "" {
		http.Error(w, "channelId is required", http.StatusBadRequest)
		return
	}

	if body.Message == nil {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	s.EmitMessageCreated(body.ChannelID, body.Message)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "message emitted",
	})
}

func (s *Server) handleEmitMessageReactionUpdated(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.ChannelID == "" {
		http.Error(w, "channelId is required", http.StatusBadRequest)
		return
	}

	if body.Message == nil {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	s.EmitMessageReactionUpdated(body.ChannelID, body.Message)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "message reaction emitted",
	})
}

func (s *Server) handleEmitMessageDeleted(w http.ResponseWriter, r *http.Request) {
	if !s.isAuthorizedInternalRequest(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body EmitMessageDeletedRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.ChannelID == "" {
		http.Error(w, "channelId is required", http.StatusBadRequest)
		return
	}

	if body.MessageID == "" {
		http.Error(w, "messageId is required", http.StatusBadRequest)
		return
	}

	s.EmitMessageDeleted(body.ChannelID, body.MessageID)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "message delete emitted",
	})
}