package log5go

import (
	"fmt"
	"sync"
)

type registry struct {
	registry map[string]Log5Go
	lock     sync.RWMutex
}

var loggerRegistry = &registry{
	make(map[string]Log5Go),
	sync.RWMutex{},
}

func (r *registry) Put(key string, logger Log5Go) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.registry[key] = logger
}

func (r *registry) Get(key string) (_ Log5Go, _ error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	logger, ok := r.registry[key]
	if !ok {
		return nil, fmt.Errorf("logger not found")
	} else {
		return logger, nil
	}
}
