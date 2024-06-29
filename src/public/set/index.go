package SgridSet

type EmptyStruct struct{}

type SgridSet[T comparable] map[T]EmptyStruct

func (s *SgridSet[T]) Add(value T) {
	(*s)[value] = EmptyStruct{}
}

func (s *SgridSet[T]) GetAll() []T {
	resu := make([]T, 0, len(*s))
	for v := range *s {
		resu = append(resu, v)
	}
	return resu
}
