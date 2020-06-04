package ml

import (
	"sync"
)

// MultipleLock is the main interface for multiLock based on key
type MultipleLock interface {
	// Lock base on the key
	Lock(interface{})

	// RLock multiLock the rw for reading
	RLock(interface{})

	// Unlock the key
	Unlock(interface{})

	// RUnlock the the read multiLock
	RUnlock(interface{})
}

func NewMultipleLock() MultipleLock {
	return &multiLock{
		locks: make(map[interface{}]*itemLock),
		lkMux: sync.Mutex{},
	}
}

type itemLock struct {
	lk  *sync.RWMutex
	cnt int64
}

// multiLock is an optimized locking system per locking key
type multiLock struct {
	locks map[interface{}]*itemLock
	lkMux sync.Mutex // protects the locks
}

func (l *multiLock) Lock(itmID interface{}) {
	l.lkMux.Lock()
	itmLock, exists := l.locks[itmID]
	if !exists {
		itmLock = &itemLock{&sync.RWMutex{}, 0}
		l.locks[itmID] = itmLock
	}
	itmLock.cnt++
	l.lkMux.Unlock()
	itmLock.lk.Lock()
}

func (l *multiLock) RLock(itmID interface{}) {
	l.lkMux.Lock()
	itmLock, exists := l.locks[itmID]
	if !exists {
		itmLock = &itemLock{&sync.RWMutex{}, 0}
		l.locks[itmID] = itmLock
	}
	itmLock.cnt++
	l.lkMux.Unlock()
	itmLock.lk.RLock()
}

func (l *multiLock) Unlock(itmID interface{}) {
	l.lkMux.Lock()
	itmLock, exists := l.locks[itmID]
	if !exists {
		panic("sync Unlock of non existent lock!!")
	}
	itmLock.lk.Unlock()
	itmLock.cnt--
	if itmLock.cnt == 0 {
		delete(l.locks, itmID)
	}
	if itmLock.cnt < 0 {
		panic("sync Unlock of free Lock!!")
	}
	l.lkMux.Unlock()
}

func (l *multiLock) RUnlock(itmID interface{}) {
	l.lkMux.Lock()
	itmLock, exists := l.locks[itmID]
	if !exists {
		panic("sync Unlock of non existent lock!!")
	}
	itmLock.lk.RUnlock()
	itmLock.cnt--
	if itmLock.cnt == 0 {
		delete(l.locks, itmID)
	}
	if itmLock.cnt < 0 {
		panic("sync Unlock of free Lock!!")
	}
	l.lkMux.Unlock()
}
