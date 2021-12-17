package websocket

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func InitManager() *Manager {
	manager := &Manager{}
	return manager
}