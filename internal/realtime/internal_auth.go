package realtime

import (
	"net/http"
	"os"
)

func (s *Server) isAuthorizedInternalRequest(r *http.Request) bool {
	expectedSecret := os.Getenv("INTERNAL_SOCKET_SECRET")
	if expectedSecret == "" {
		return true
	}

	receivedSecret := r.Header.Get("X-Internal-Secret")

	return receivedSecret == expectedSecret
}