package log4go

import (
  "fmt"
  "sync"
)

type registry struct {
  registry map[string]Log4Go
  lock sync.RWMutex
}

var loggerRegistry = &registry{
  make(map[string]Log4Go),
  sync.RWMutex{},
}

func (r *registry) Put(key string, logger Log4Go) error {
  r.lock.Lock()
  defer r.lock.Unlock()

  fmt.Printf("putting logger for key: %s\n", key)

  if _, ok := r.registry[key]; ok {
    return fmt.Errorf("logger already exists for key: %s", key)
  }

  r.registry[key] = logger
  return nil
}

func (r *registry) Get(key string) (_ Log4Go, _ error) {
  r.lock.RLock()
  defer r.lock.RUnlock()

  logger, ok := r.registry[key]
  if !ok {
    return nil, fmt.Errorf("logger not found")
  } else {
    return logger, nil
  }
}
