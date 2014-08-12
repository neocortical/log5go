package log5go

import (
  "fmt"
  "sync"
)

type registry struct {
  registry map[string]log5go
  lock sync.RWMutex
}

var loggerRegistry = &registry{
  make(map[string]log5go),
  sync.RWMutex{},
}

func (r *registry) Put(key string, logger log5go) error {
  r.lock.Lock()
  defer r.lock.Unlock()

  if _, ok := r.registry[key]; ok {
    return fmt.Errorf("logger already exists for key: %s", key)
  }

  r.registry[key] = logger
  return nil
}

func (r *registry) Get(key string) (_ log5go, _ error) {
  r.lock.RLock()
  defer r.lock.RUnlock()

  logger, ok := r.registry[key]
  if !ok {
    return nil, fmt.Errorf("logger not found")
  } else {
    return logger, nil
  }
}
