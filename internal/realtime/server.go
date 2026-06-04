package realtime

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/beyremagrebi/combo_app_socket/internal/config"
	"github.com/zishang520/socket.io/servers/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

type Server struct {
	config config.Config
	io     *socket.Server
	mux    *http.ServeMux
}

func NewServer(cfg config.Config) *Server {
	server := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}

	server.io = server.createSocketServer()
	server.registerHandlers()
	server.registerRoutes()

	return server
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort)

	log.Printf("Socket.IO Go server running on http://localhost:%s", s.config.AppPort)
	log.Println("Socket.IO path registered on /socket.io/")

	return http.ListenAndServe(address, s.mux)
}

func (s *Server) registerRoutes() {
	s.mux.Handle("/socket.io/", s.io.ServeHandler(nil))

	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Socket server is running"))
	})

	s.registerInternalRoutes()
	s.registerInternalMessageRoutes()
	s.registerInternalDirectMessageRoutes()
}

func (s *Server) createSocketServer() *socket.Server {
	opts := socket.DefaultServerOptions()

	opts.SetPingTimeout(20 * time.Second)
	opts.SetPingInterval(25 * time.Second)

	opts.SetCors(&types.Cors{
		Origin:      s.config.ClientOrigin,
		Credentials: true,
	})

	return socket.NewServer(nil, opts)
}