// 链式操作 string
package replace

import (
	"fmt"
	"strings"
)

const (
	WHERE      = "${WHERE}"
	SELECTS    = "${SELECTS}"
	PAGINATION = "${PAGINATION}"
)

type ReplaceChain struct {
	originString  string
	currentString string
	isDebug       bool
	count         *int64
}

func BuildReplaceChain(origin string) *ReplaceChain {
	return &ReplaceChain{
		originString:  origin,
		currentString: origin,
		isDebug:       false,
		count:         new(int64),
	}
}

func (r *ReplaceChain) GetOriginString() string {
	return r.originString
}

func (r *ReplaceChain) Debug() *ReplaceChain {
	r.isDebug = true
	return r
}

func (r *ReplaceChain) Replace(old, new string) *ReplaceChain {
	r.currentString = strings.Replace(r.currentString, old, new, 1)
	if r.isDebug {
		fmt.Println("r.currentString \n", r.currentString)
	}
	return r
}

func (r *ReplaceChain) Get() string {
	return r.currentString
}

func (r *ReplaceChain) Reset() *ReplaceChain {
	r.currentString = r.originString
	return r
}

// 数据库操作
func (r *ReplaceChain) ReplaceWhere(new string) *ReplaceChain {
	return r.Replace(WHERE, new)
}

func (r *ReplaceChain) ReplaceSelects(new string) *ReplaceChain {
	return r.Replace(SELECTS, new)
}

func (r *ReplaceChain) ReplacePagination(offset int, size int) *ReplaceChain {
	if size == 0 {
		size = 10
	}
	return r.Replace(PAGINATION, fmt.Sprintf(" limit %v,%v", offset, size))
}

func (r *ReplaceChain) ReplaceAsCount() *ReplaceChain {
	return r.Replace(SELECTS, " count(*) as total ")
}

func (r *ReplaceChain) ReplaceWithNoPagination() *ReplaceChain {
	return r.Replace(PAGINATION, "")
}

func (r *ReplaceChain) GetCountVo() *int64 {
	return r.count
}
