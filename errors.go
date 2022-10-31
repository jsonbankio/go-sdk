package jsonbank

import "fmt"

type RequestError struct {
	Code    string
	Message string
}

func (error *RequestError) Error() string {
	return fmt.Sprintf("[%v]: %v", error.Code, error.Message)
}
