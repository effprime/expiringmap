package expiringmap

import (
	"fmt"
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
	PanicFull bool
}

func (s *Settings) Default() {
	if s.Age == 0 {
		s.Age = AgeDefault
	}

	if s.MaxLength == 0 {
		s.MaxLength = LengthDefault
	}
}

func NewExpiringMap[T any](s Settings) *ExpiringMap[T] {
	s.Default()
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

	if len(e.m) >= e.config.MaxLength {
		if e.config.PanicFull {
			panic(fmt.Sprintf("expiring map is full (max length: %v)", e.config.MaxLength))
		}
		delete(e.m, e.oldestKey())
	}
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

func (e *ExpiringMap[T]) oldestKey() string {
	oldestKey := ""
	oldestValue := time.Time{}
	for k, v := range e.m {
		if oldestValue.IsZero() || v.Timestamp.Before(oldestValue) {
			oldestKey = k
			oldestValue = v.Timestamp
		}
	}
	return oldestKey
}
