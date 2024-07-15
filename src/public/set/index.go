package SgridSet

import "sync"

type ES struct{}

type SgridSet[T comparable] struct {
	v  map[T]*ES     // value
	c  int           // count
	rw *sync.RWMutex // rw lock
}

func NewSgridSet[T comparable]() *SgridSet[T] {
	return &SgridSet[T]{
		rw: &sync.RWMutex{},
		c:  0,
		v:  make(map[T]*ES, 1024),
	}
}

func (s *SgridSet[T]) Add(value ...T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, V := range value {
		if s.v[V] == nil {
			s.c = s.c + 1
		}
		s.v[V] = &ES{}
	}
}

func (s *SgridSet[T]) GetAll() []T {
	s.rw.RLock()
	defer s.rw.RUnlock()
	resu := make([]T, 0, len(s.v))
	for v := range s.v {
		resu = append(resu, v)
	}
	return resu
}

func (s *SgridSet[T]) Remove(value T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	if s.v[value] != nil {
		s.c = s.c - 1
	}
	delete(s.v, value)
}

func (s *SgridSet[T]) GetCount() int {
	return s.c
}
