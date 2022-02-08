package websockets

import (
	"net/http"
)

type Server struct {
	Manager  *Manager
	Protocol Protocol
	Subscriber Subscriber
}

func InitServer(protocol Protocol, Subscriber Subscriber) *Server {
	return &Server{
		Manager:  InitManager(),
		Protocol: protocol,
		Subscriber: Subscriber,
	}
}

func Listen(addr string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
