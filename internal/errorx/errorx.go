package errorx

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrNotFound generic not found error
	ErrNotFound = errors.New("not found")
	// ErrInvalidArgument generic invalid argument error
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrUnknown generic unknown error
	ErrUnknown = errors.New("unknown")
	// ErrConfig configuration error
	ErrConfig = errors.New("configuration error")
	// ErrAlreadyExists already exists error
	ErrAlreadyExists = errors.New("already exists")
	// ErrEOF error
	ErrEOF = fmt.Errorf("EOF: %w", io.EOF)
	// ErrOutOfRange index out of range error
	ErrOutOfRange = errors.New("out of range")
)
