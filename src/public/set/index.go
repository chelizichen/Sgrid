// 空结构体本身确实会占用一些内存，但通常这种占用是非常小的。在Go中，即使结构体没有字段，它仍然会被分配一个非零大小的内存空间。这是因为每个对象在内存中都需要一个唯一的地址，即使是空结构体也不例外。然而，这个内存占用通常可以忽略不计，因为它远远小于一个指针的大小。
// 空结构体主要用于以下几种情况：
// 作为占位符：当你需要一个类型的占位符，但不需要存储任何数据时，可以使用空结构体。例如，在通道（channel）中传递信号或事件时，你可能不需要传递任何具体的数据，只需要知道某个事件发生了。
// 作为方法的接收器：当你想为一个不包含任何数据的类型定义方法时，可以使用空结构体。这种方法常见于只关注行为而不关注状态的场景。
// 节省空间：尽管空结构体本身仍会占用一些内存，但相比于包含多个字段的结构体，它的内存占用是极小的。如果你正在处理大量数据且每个数据项都包含这种结构体，使用空结构体可以帮助节省内存。
// 作为映射的键：有时，你可能需要一个唯一的标识符作为映射的键，但不需要存储与该键关联的任何数据。在这种情况下，可以使用空结构体作为键的类型。
// 需要注意的是，虽然空结构体本身占用内存很小，但如果你在切片、映射或通道中大量使用它，这些集合类型本身会占用更多的内存来管理它们的内部结构（如长度、容量、指针等）。
// 总的来说，空结构体在Go中是一个有用的工具，尽管它本身会占用一些内存，但在某些情况下使用它可以带来代码上的清晰和灵活性。
// 如下，可以用来作为一个简单的Set

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