package server

import (
	"easycast/server/context"
	"sync"
)

type CnMap struct {
	lock sync.Mutex
	cmap map[*context.WsContext]interface{}
}

func NewConnectionStorage() *CnMap {
	m := new(CnMap)
	m.cmap = make(map[*context.WsContext]interface{})
	return m
}

func (m *CnMap) GetAndLock() map[*context.WsContext]interface{} {
	m.lock.Lock()
	return m.cmap
}

func (m *CnMap) UnLock() {
	m.lock.Unlock()
}

func (m *CnMap) Put(c *context.WsContext) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cmap[c] = nil
}

func (m *CnMap) TryRemove(c *context.WsContext) {
	// TODO remove ctx
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.cmap, c)
}
