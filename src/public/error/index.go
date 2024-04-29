package sgridError

import (
	"errors"
	"fmt"
)

var Reflect_Get_Method_Error = func(m string) error { return New(-10, fmt.Sprintf("Reflect_Get_Method_Error :: %v", m)) }
var Reflect_Get_Field_Error = func(f string) error { return New(-11, fmt.Sprintf("Reflect_Get_Field_Error :: %v", f)) }
var Request_Error = func(f string) error { return New(-11, fmt.Sprintf("Request_Error :: %v", f)) }

func New(code int, msg string) (err error) {
	errorInfo := fmt.Sprintf("sgrid/error code :: %v \n message %v", code, msg)
	return errors.New(errorInfo)
}
