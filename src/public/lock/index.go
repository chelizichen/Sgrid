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
