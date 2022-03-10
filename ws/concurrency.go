package ws

import (
	"sync"

	"bgtools-api/models"
)

type ConnMap struct {
	m  map[string]*WsConnection
	mu sync.Mutex
}

type RoomMap struct {
	m  map[string]models.RoomInfoSet
	mu sync.RWMutex
}

type PlayerMap struct {
	m  map[string]string
	mu sync.RWMutex
}

func NewConnMap() *ConnMap {
	return &ConnMap{
		m: make(map[string]*WsConnection),
	}
}

func NewRoomMap() *RoomMap {
	return &RoomMap{
		m: make(map[string]models.RoomInfoSet),
	}
}

func NewPlayerMap() *PlayerMap {
	return &PlayerMap{
		m: make(map[string]string),
	}
}

func (c *ConnMap) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.m)
}

func (c *ConnMap) Get(id string) (*WsConnection, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.m[id]
	return v, ok
}

func (c *ConnMap) Set(id string, conn *WsConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[id] = conn
}

func (c *ConnMap) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.m, id)
}

func (c *ConnMap) GetKeys() map[string]struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	res := make(map[string]struct{}, len(c.m))

	for k := range c.m {
		res[k] = struct{}{}
	}

	return res
}

func (r *RoomMap) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.m)
}

func (r *RoomMap) Get(id string) (models.RoomInfoSet, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.m[id]
	return v, ok
}

func (r *RoomMap) Set(id string, room models.RoomInfoSet) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.m[id] = room
}

func (r *RoomMap) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.m, id)
}

func (r *RoomMap) Range(f func(id string, room models.RoomInfoSet)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for k, v := range r.m {
		f(k, v)
	}
}

func (p *PlayerMap) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.m)
}

func (p *PlayerMap) Get(id string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	v, ok := p.m[id]
	return v, ok
}

func (p *PlayerMap) Set(connid string, roomid string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.m[connid] = roomid
}

func (p *PlayerMap) Delete(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.m, id)
}
