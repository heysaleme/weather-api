package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrUpstream     = errors.New("upstream error")
)

func InvalidInput(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, message)
}

func NotFound(message string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, message)
}

func Upstream(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrUpstream, fmt.Sprintf(format, args...))
}
