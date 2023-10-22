package telegram

import (
	"fmt"
)

type RequestError struct {
	err error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

func NewRequestError(e error) error {
	return &RequestError{
		err: fmt.Errorf("request error: %w", e),
	}
}
