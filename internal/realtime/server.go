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
}

func NewServer(cfg config.Config) *Server {
	server := &Server{
		config: cfg,
	}

	server.io = server.createSocketServer()
	server.registerHandlers()

	return server
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort)

	http.Handle("/socket.io/", s.io.ServeHandler(nil))

	log.Printf("Socket.IO Go server running on http://localhost:%s", s.config.AppPort)

	return http.ListenAndServe(address, nil)
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