package task

import (
	"sync"
)

type Lock struct {
	mu sync.RWMutex
}

func (l *Lock) RLock() func() {
	l.mu.RLock()
	return func() { l.mu.RUnlock() }
}
func (l *Lock) WLock() func() {
	l.mu.Lock()
	return func() { l.mu.Unlock() }
}
