package sky

import "errors"

var (
	ErrInternal     = errors.New("internal error")
	ErrInvalidJson  = errors.New("invalid json")
	ErrPersonValid  = errors.New("invalid person data")
	ErrNoneAffected = errors.New("operation has no effect")
)
