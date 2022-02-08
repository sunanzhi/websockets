package websockets

import (
	cmap "github.com/orcaman/concurrent-map"
)

type Manager struct {
	Clients cmap.ConcurrentMap
}

func InitManager() *Manager {
	return &Manager{
		Clients: cmap.New(),
	}
}

// Get 获取客户端
func (m *Manager) Get(id string) *Client {
	if client, ok := m.Clients.Get(id); ok {
		return client.(*Client)
	}

	return nil
}

// Join 加入
func (m *Manager) Join(c *Client) bool {
	m.Clients.Set(c.Id, c)

	return true
}

// Leave 离开
func (m *Manager) Leave(c *Client) bool {
	m.Clients.Remove(c.Id)

	return true
}

// Broadcast 广播
func (m *Manager) Broadcast(msg SubscriberBody) {
	go func() {
		for v := range m.Clients.IterBuffered() {
			client := v.Val.(*Client)
			if client.IsServer {
				go client.Server.Protocol.Send(client, msg)
			}
		}
	}()
}