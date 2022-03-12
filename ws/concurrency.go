package ws

import (
	"sync"

	"bgtools-api/models"
)

// <summary>: プレイヤー情報向けスレッドセーフなデータ格納庫
type PlayerMap struct {
	m  map[string]PlayerConn
	mu sync.Mutex
}

// <summary>: 部屋情報向けスレッドセーフなデータ格納庫
type RoomMap struct {
	m  map[string]models.RoomInfoSet
	mu sync.RWMutex
}

// <summary>: プレイヤー情報格納庫の初期化
func NewPlayerMap() *PlayerMap {
	return &PlayerMap{
		m: make(map[string]PlayerConn),
	}
}

// <summary>: 部屋情報格納庫の初期化
func NewRoomMap() *RoomMap {
	return &RoomMap{
		m: make(map[string]models.RoomInfoSet),
	}
}

// <summary>: プレイヤーマップにあるデータの数を数えます
func (p *PlayerMap) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.m)
}

// <summary>: プレイヤーマップからキーをもとに情報を取得します
func (p *PlayerMap) Get(id string) (PlayerConn, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.m[id]
	return v, ok
}

// <summary>: プレイヤーマップに情報を格納します
func (p *PlayerMap) Set(id string, conn PlayerConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.m[id] = conn
}

// <summary>: プレイヤーマップに部屋情報を上書きします
// <remark>: 上書きの成否が取得できます
func (p *PlayerMap) SetRoomId(connid, roomid string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.m[connid]

	if ok {
		v.RoomId = roomid
		p.m[connid] = v
	}

	return ok
}

// <summary>: プレイヤーマップから情報を削除します
func (p *PlayerMap) Delete(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.m, id)
}

// <summary>: プレイヤーマップからプレイヤーがいる部屋情報の一覧を取得します
func (p *PlayerMap) PlayerRoomData() map[string]string {
	p.mu.Lock()
	defer p.mu.Unlock()

	result := make(map[string]string, len(p.m))

	for k, v := range p.m {
		result[k] = v.RoomId
	}

	return result
}

// <summary>: 部屋マップにあるデータの数を数えます
func (r *RoomMap) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.m)
}

// <summary>: 部屋マップからキーをもとに情報を取得します
func (r *RoomMap) Get(id string) (models.RoomInfoSet, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.m[id]
	return v, ok
}

// <summary>: 部屋マップに情報を格納します
func (r *RoomMap) Set(id string, room models.RoomInfoSet) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.m[id] = room
}

// <summary>: 部屋マップから情報を削除します
func (r *RoomMap) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.m, id)
}

// <summary>: 部屋マップの情報に対して、一連の処理を実行します
func (r *RoomMap) Range(f func(id string, room models.RoomInfoSet)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for k, v := range r.m {
		f(k, v)
	}
}
