package common

import "sync"

var (
	idmutex sync.Mutex
	id      uint64 = 1000
)

const UINT64_MAX uint64 = ^uint64(0)

func GenerateID() uint64 {
	idmutex.Lock()
	defer idmutex.Unlock()
	if id == UINT64_MAX {
		id = 1000
	}
	id++
	return id
}
