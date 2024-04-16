package sgridLock

import (
	"runtime"
	"sync/atomic"
)

const MaxLockFaileTime = 36

type SgridLock struct {
	targetLock atomic.Bool
	backoff    int8
}

func (s *SgridLock) Lock() bool {
	b := s.targetLock.Swap(true)
	if b {
		if s.backoff <= MaxLockFaileTime {
			s.backoff <<= 1
		}
		defer func() {
			for i := 0; i < int(s.backoff); i++ {
				runtime.Gosched()
			}
		}()
		return false
	}
	s.backoff = 0
	return true
}

func (s *SgridLock) Unlock() {
	s.targetLock.Store(false)
}

const SyncSateMapMaxFaileTime = 4

type SyncSafetyMap[K comparable, V comparable] struct {
	value      map[K]V
	targetLock atomic.Bool
	backoff    int8
}

func (s *SyncSafetyMap[K, V]) Set(key K, Value V) bool {
	return s.lock(func() {
		s.value[key] = Value
	})
}

func (s *SyncSafetyMap[K, V]) Delete(key K) bool {
	return s.lock(func() {
		delete(s.value, key)
	})
}

func (s *SyncSafetyMap[K, V]) Get(key K) V {
	return s.value[key]
}

func (s *SyncSafetyMap[K, V]) lock(event func()) bool {
	b := s.targetLock.Swap(true)
	if b {
		if s.backoff <= SyncSateMapMaxFaileTime {
			s.backoff <<= 1
		}
		defer func() {
			for i := 0; i < int(s.backoff); i++ {
				runtime.Gosched()
			}
		}()
		return false
	}
	defer s.unlock()
	event()
	s.backoff = 0
	return true
}

func (s *SyncSafetyMap[K, V]) unlock() {
	s.targetLock.Store(false)
}
