package expiringmap

import (
	"sync"
	"time"
)

const (
	AgeDefault    = 5 * time.Minute
	LengthDefault = 1000
)

type ExpiringMap[T any] struct {
	lock   sync.Mutex
	config Settings
	m      map[string]DatedValue[T]
}

type Settings struct {
	Age       time.Duration
	MaxLength int
}

func (s *Settings) Default() {
	s.Age = AgeDefault
	s.MaxLength = LengthDefault
}

func NewExpiringMap[T any](s Settings) *ExpiringMap[T] {
	m := map[string]DatedValue[T]{}
	return &ExpiringMap[T]{
		config: s,
		m:      m,
	}
}

type DatedValue[T any] struct {
	Value     T
	Timestamp time.Time
}

func (e *ExpiringMap[T]) Set(key string, value T) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.m[key] = DatedValue[T]{Value: value, Timestamp: time.Now()}
}

func (e *ExpiringMap[T]) Remove(key string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.m, key)
}

func (e *ExpiringMap[T]) Get(key string) (T, bool) {
	e.lock.Lock()
	defer e.lock.Unlock()

	val, ok := e.m[key]
	if !ok {
		return val.Value, false
	}
	if time.Now().Sub(val.Timestamp) > e.config.Age {
		delete(e.m, key)
		val, _ := e.m[key]
		return val.Value, false
	}
	return val.Value, true
}
