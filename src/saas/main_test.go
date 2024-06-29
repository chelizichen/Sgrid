package saas_test

import (
	"Sgrid/src/saas"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAuthString(t *testing.T) {
	saas.UsesaasPerm.StartTime = "2024-01-01 00:00:00"
	saas.UsesaasPerm.EndTime = "2024-06-28 00:00:00"
	saas.UsesaasPerm.Username = "chelizichen"
	saas.UsesaasPerm.Password = "leemulus21"
	s, err := saas.UsesaasPerm.GenAuthString()
	fmt.Println("s", s)
	assert.Nil(t, err)

}

func TestParseAuthString(t *testing.T) {
	st := `QUFBQUFBQUFBQUFBQUFBQUFBQUFBSUVSbFJvUlJXbGFqSU1XYkMyY1NxdHdMVGJ6UC81Q1VyazRLcVdCSGdZQk81dTlBenlPWDlsREZ5UnQ4TmV5RXliN2hxdDM3MFczY2ZVeEdYUGw=`
	saas.UsesaasPerm.SetAuthString(st)
	s, err := saas.UsesaasPerm.ParseAuthString()
	assert.Nil(t, err)
	fmt.Println(s)
	arguments := strings.Split(s, "\n")
	if len(arguments) != 4 {
		panic("error.length not equal 4")
	}
	saas.UsesaasPerm.Username = arguments[0]
	saas.UsesaasPerm.Password = arguments[1]
	saas.UsesaasPerm.StartTime = arguments[2]
	saas.UsesaasPerm.EndTime = arguments[3]
	fmt.Println("sass.UsesassPerm", saas.UsesaasPerm)
}

func neverRun() {
	var expireTestString = "QUFBQUFBQUFBQUFBQUFBQUFBQUFBSUVSbFJvUlJXbGFqSU1XYkMyY1NxdHdMVGJ6UC81Q1VyazRLcVdCSGdZQk81dTlBenlPWDlsREZ5UnQ4TmV6RXliOGhxbCs3MFczY2ZVeEdYUGw="
	fmt.Println(expireTestString)
	var noExpireTimeStr = `QUFBQUFBQUFBQUFBQUFBQUFBQUFBSUVSbFJvUlJXbGFqSU1XYkMyY1NxdHdMVGJ6UC81Q1VyazRLcVdCSGdZQk81dTlBenlPWDlsREZ5UnQ4TmV5RXliN2hxdDM3MFczY2ZVeEdYUGw=`
	fmt.Println(noExpireTimeStr)
}
