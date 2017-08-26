package server

import (
	"sync"
	"net"
)

type CnMap struct {
	lock sync.Mutex
	cmap map[string]net.Conn
}

func NewConnectionStorage() *CnMap {
	m := new(CnMap)
	m.cmap = make(map[string]net.Conn)
	return m
}

func (m *CnMap) GetAndLock() map[string]net.Conn {
	m.lock.Lock()
	return m.cmap
}

func (m *CnMap) UnLock() {
	m.lock.Unlock()
}

func (m *CnMap) Put(c string, conn net.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cmap[c] = conn
}

func (m *CnMap) TryRemove(c string) {
	// TODO remove ctx
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.cmap, c)
}
