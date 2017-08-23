package server

import (
	"sync"

	"github.com/shmel1k/exchangego/exchange/session/context"
)

type CnMap struct {
	lock sync.Mutex
	cmap map[*context.ExContext]interface{}
}

func NewConnectionStorage() *CnMap {
	m := new(CnMap)
	m.cmap = make(map[*context.ExContext]interface{})
	return m
}

func (m *CnMap) GetAndLock() map[*context.ExContext]interface{} {
	m.lock.Lock()
	return m.cmap
}

func (m *CnMap) UnLock() {
	m.lock.Unlock()
}

func (m *CnMap) Put(c *context.ExContext) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cmap[c] = nil
}

func (m *CnMap) TryRemove(c *context.ExContext) {
	// TODO remove ctx
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.cmap, c)
}
