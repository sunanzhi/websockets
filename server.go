package websocket

import "net"

type Server struct {
	manager      *Manager
	listener     net.Listener
}

func NewServer(listener net.Listener) *Server {
	return &Server{
		manager:      InitManager(),
		listener:     listener,
	}
}
