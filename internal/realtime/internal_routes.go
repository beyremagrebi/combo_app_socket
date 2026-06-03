package realtime

import (
	"encoding/json"
	"net/http"
	"os"
)

type EmitAuditLogRequest struct {
	TeamID string         `json:"teamId"`
	Log    map[string]any `json:"log"`
}

func (s *Server) registerInternalRoutes() {
	s.mux.HandleFunc("/internal/audit-logs/emit", s.handleEmitAuditLog)
}

func (s *Server) handleEmitAuditLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedSecret := os.Getenv("INTERNAL_SOCKET_SECRET")
	receivedSecret := r.Header.Get("X-Internal-Secret")

	if expectedSecret != "" && receivedSecret != expectedSecret {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body EmitAuditLogRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.TeamID == "" {
		http.Error(w, "teamId is required", http.StatusBadRequest)
		return
	}

	if body.Log == nil {
		http.Error(w, "log is required", http.StatusBadRequest)
		return
	}

	s.EmitAuditLogCreated(body.TeamID, body.Log)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "audit log emitted",
	})
}