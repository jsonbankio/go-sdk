package jsonbank

import "fmt"

type RequestError struct {
	Code    string
	Message string
}

func (error *RequestError) Error() string {
	return fmt.Sprintf("[%v]: %v", error.Code, error.Message)
}

var InvalidJsonError = RequestError{"invalid_json_content", "Content is not a valid JSON string"}
