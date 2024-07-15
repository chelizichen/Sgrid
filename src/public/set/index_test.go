package SgridSet_test

import (
	set "Sgrid/src/public/set"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BASE_SET(t *testing.T) {
	s := set.NewSgridSet[string]()
	s.Add("11")
	s.Add("22")
	s.Add("33")
	threeElement := s.GetAll()
	assert.Len(t, threeElement, 3)
	assert.EqualValues(t, s.GetCount(), 3)

	s.Add("33")
	alsoThreeElement := s.GetAll()
	assert.Len(t, alsoThreeElement, 3)
	assert.EqualValues(t, s.GetCount(), 3)

	s.Remove("11")
	twoElement := s.GetAll()
	assert.Len(t, twoElement, 2)
	assert.EqualValues(t, s.GetCount(), 2)
}
