package SgridHashUtil

type IHash interface {
	Sum64(string) uint64
}

func NewDefaultHash() IHash {
	return fnv64a{}
}

type fnv64a struct{}

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// Sum64 gets the string and returns its uint64 hash value.
func (f fnv64a) Sum64(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}
