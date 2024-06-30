package SgridSet

type ES struct{}

type SgridSet[T comparable] map[T]ES

func (s *SgridSet[T]) Add(value T) {
	(*s)[value] = ES{}
}

func (s *SgridSet[T]) GetAll() []T {
	resu := make([]T, 0, len(*s))
	for v := range *s {
		resu = append(resu, v)
	}
	return resu
}
