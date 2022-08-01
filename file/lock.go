package file

import (
	"sync"
)

type Lock struct{ mu sync.RWMutex }

func (lock *Lock) RLock() func() {
	lock.mu.RLock()
	return func() { lock.mu.RUnlock() }
}
func (lock *Lock) Lock() func() {
	lock.mu.Lock()
	return func() { lock.mu.Unlock() }
}
